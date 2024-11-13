from collections import defaultdict
import os
import shutil
import subprocess
from typing import Dict, List
from dateutil.parser import isoparse
import json
import time
import socket
import re


# import os
import sys
from pathlib import Path
import tempfile
import requests
import toml
from pprint import pprint
import urllib.parse
import asyncio

DEBUG = True
DEFAULT_DENOM = "agnet"
BASE_DENOM = DEFAULT_DENOM
DISPLAY_DENOM = "gnet"
DEFAULT_CHAIN_BINARY = "../../build/galacticad"

DEFAULT_TEST_MONIKER = "test-node01"
DEFAULT_TEST_CHAINID = "test_41239-41239"

PREDEFINED_KEY_MNEMONIC_TREASURY = "gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"
PREDEFINED_KEY_MNEMONIC_FAUCET = "heart grape ignore face equip monkey keep armor tumble donkey final horror harsh way retire this enforce pave there unfair scrap shine physical since"
PREDEFINED_KEY_MNEMONIC_NODE_KEYS = [
    "kick treat protect present permit business own nuclear ranch ancient around deposit dignity cabin kiwi parade sister market must crime tag update yellow theory",
    "minimum sing arrow way comfort obvious purse piece reward simple fitness fence october dutch genius spike sunset empower limit dog dutch kid online file",
    "uniform spread february wife quality device mix fish rapid win improve van eagle target icon home charge birth reward slogan season robust thunder over",
]


class GnetAmount:
    DENOMINATIONS = {"agnet": 1, "ugnet": 1e12, "mgnet": 1e15, "gnet": 1e18}

    def __init__(self, amount):
        if isinstance(amount, float):
            self.amount = amount
        else:
            match = re.match(r"([\d\.]+)(agnet|ugnet|mgnet|gnet)", amount)
            if match:
                self.amount = float(match.group(1)) * self.DENOMINATIONS[match.group(2)]
            else:
                raise ValueError(f"Invalid amount format: {amount}")

    def __repr__(self):
        for denomination, value in sorted(
            self.DENOMINATIONS.items(), key=lambda item: -item[1]
        ):
            if self.amount % value == 0:
                return f"{self.__class__.__name__}({int(self.amount // value)}{denomination})"

    def __str__(self):
        for denomination, value in sorted(
            self.DENOMINATIONS.items(), key=lambda item: item[1]
        ):
            if self.amount % value == 0:
                return f"{str(int(self.amount // value))}{denomination}"

    def __format__(self, format_spec):
        return str(self)

    def __add__(self, other):
        if not isinstance(other, GnetAmount):
            other = GnetAmount(other)
        return GnetAmount(self.amount + other.amount)
        # else:
        #     raise TypeError(
        #         f"unsupported operand type(s) for +: '{self.__class__.__name__}' and '{type(other).__name__}'"
        #     )

    def __sub__(self, other):
        if isinstance(other, GnetAmount):
            return GnetAmount(self.amount - other.amount)
        else:
            raise TypeError(
                f"unsupported operand type(s) for -: '{self.__class__.__name__}' and '{type(other).__name__}'"
            )

    def __eq__(self, other):
        if isinstance(other, GnetAmount):
            return self.amount == other.amount
        else:
            return False

    def __lt__(self, other):
        if isinstance(other, GnetAmount):
            return self.amount < other.amount
        else:
            raise TypeError(
                f"unsupported operand type(s) for <: '{self.__class__.__name__}' and '{type(other).__name__}'"
            )

    def __mul__(self, other):
        if isinstance(other, (int, float)):
            return GnetAmount(self.amount * other)
        else:
            raise TypeError(
                f"unsupported operand type(s) for *: '{self.__class__.__name__}' and '{type(other).__name__}'"
            )

    def __truediv__(self, other):
        if isinstance(other, (int, float)):
            return GnetAmount(self.amount / other)
        elif isinstance(other, (GnetAmount)):
            return float(self.amount / other.amount)
        else:
            raise TypeError(
                f"unsupported operand type(s) for /: '{self.__class__.__name__}' and '{type(other).__name__}'"
            )

    def min_denom_amount_str(self) -> str:
        return str(int(self.amount))


def wait_for_port(port, host="127.0.0.1", timeout=40.0):
    start_time = time.perf_counter()
    while True:
        try:
            with socket.create_connection((host, port), timeout=timeout):
                break
        except OSError as ex:
            time.sleep(0.1)
            if time.perf_counter() - start_time >= timeout:
                raise TimeoutError(
                    "Waited too long for the port {} on host {} to start accepting "
                    "connections.".format(port, host)
                ) from ex


def get_current_height(cli):
    try:
        status = cli.status()
    except AssertionError as e:
        print(f"get sync status failed: {e}", file=sys.stderr)
    else:
        current_height = int(status["sync_info"]["latest_block_height"])
    return current_height


def wait_for_block(cli, height, timeout=240):
    for _ in range(timeout * 2):
        current_height = get_current_height(cli)
        if current_height >= height:
            break
        print("current block height", current_height)
        time.sleep(0.5)
    else:
        raise TimeoutError(f"wait for block {height} timeout")


def interact(cmd, ignore_error=False, input=None, **kwargs):
    if DEBUG:
        print("\033[94m" + cmd + "\033[0m")
    kwargs.setdefault("stderr", subprocess.STDOUT)
    proc = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        shell=True,
        **kwargs,
    )
    # begin = time.perf_counter()
    stdout, _ = proc.communicate(input=input)
    # print('[%.02f] %s' % (time.perf_counter() - begin, cmd))
    if not ignore_error:
        assert proc.returncode == 0, f'{stdout.decode("utf-8")} ({cmd})'
    return stdout


def safe_cli_string(s):
    'wrap string in "", used for cli argument when contains spaces'
    if len(f"{s}".split()) > 1:
        return f"'{s}'"
    return f"{s}"


def build_cli_args_safe(*args, **kwargs):
    args = [safe_cli_string(arg) for arg in args if arg]
    for k, v in kwargs.items():
        if v is None:
            continue
        args.append("--" + k.strip("_").replace("_", "-"))
        args.append(safe_cli_string(v))
    return list(map(str, args))


def build_cli_args(*args, **kwargs):
    args = [arg for arg in args if arg is not None]
    for k, v in kwargs.items():
        if v is None:
            continue
        args.append("--" + k.strip("_").replace("_", "-"))
        args.append(v)
    return list(map(str, args))


def format_doc_string(**kwargs):
    def decorator(target):
        target.__doc__ = target.__doc__.format(**kwargs)
        return target

    return decorator


class GalaToml:
    def __init__(self, path):
        self.path = path
        self.load()

    def load(self):
        "Load config from toml file at self.path"
        try:
            with open(self.path, "r") as file:
                self.config = toml.load(file)
        except FileNotFoundError:
            print(f"Config file {self.path} not found.")

    def save(self):
        "Save config to self.path in toml format"
        with open(self.path, "w") as file:
            toml.dump(self.config, file)

    def deep_update(self, original, new):
        "Recursive update of nested dictionaries"
        for key, value in new.items():
            if isinstance(value, dict) and isinstance(original.get(key), dict):
                self.deep_update(original[key], value)
            else:
                original[key] = value

    def edit(self, new_config):
        "Edit config with new_config, preserving unchanged keys"
        self.deep_update(self.config, new_config)
        self.save()

    def diff(self, other):
        "Return a dictionary containing the differences between self and other"
        diff_config = {}
        for key, value in self.config.items():
            if key not in other.config:
                diff_config[key] = value
            elif isinstance(value, dict):
                sub_diff = self._diff_dict(value, other.config[key])
                if sub_diff:
                    diff_config[key] = sub_diff
            elif value != other.config[key]:
                diff_config[key] = value
        for key, value in other.config.items():
            if key not in self.config:
                diff_config[key] = value
        return diff_config

    def _diff_dict(self, dict1, dict2):
        "Helper method to compare two dictionaries"
        diff_dict = {}
        for key, value in dict1.items():
            if key not in dict2:
                diff_dict[key] = value
            elif isinstance(value, dict):
                sub_diff = self._diff_dict(value, dict2[key])
                if sub_diff:
                    diff_dict[key] = sub_diff
            elif value != dict2[key]:
                diff_dict[key] = value
        return diff_dict

    def apply_addr(self, new_addr):
        "Replace the host in all address values in config with the host from new_addr"
        new_host = urllib.parse.urlparse(new_addr).hostname or new_addr
        self._apply_addr_to_dict(self.config, new_host)
        self.save()

    def _apply_addr_to_dict(self, inner_dict, new_host):
        "Helper method to apply new_host to nested dictionaries recursively"
        address_suffixes = ["address", "addr", "proxy_app"]

        for key, value in inner_dict.items():
            if isinstance(value, dict):
                self._apply_addr_to_dict(value, new_host)
            elif any(
                key.endswith(suffix) for suffix in address_suffixes
            ) and isinstance(value, str):
                # TODO: make sure that addresses with '' are not in conflict
                if value == "":
                    continue
                parsed_url = urllib.parse.urlparse(value)
                scheme = parsed_url.scheme
                invalid_scheme = not scheme or scheme == "localhost"
                if invalid_scheme:
                    try:
                        hostname, port = value.split(":")
                    except ValueError as e:
                        print(e, value)
                else:
                    port = parsed_url.port
                netloc = new_host + ":" + str(port)
                if invalid_scheme:
                    new_endpoint = netloc
                else:
                    new_endpoint = urllib.parse.urlunparse(
                        (scheme, netloc, *parsed_url[2:])
                    )
                inner_dict[key] = new_endpoint


class GalaClientConfig(GalaToml):
    def __init__(self, path):
        super().__init__(path)

    def to_dict(self):
        return self.config


class BinaryCommand:
    def __init__(self, cmd):
        self.cmd = cmd

    def __call__(self, cmd, *args, stdin=None, **kwargs):
        "execute cmd with binary chaind"
        args = " ".join(build_cli_args_safe(cmd, *args, **kwargs))
        return interact(f"{self.cmd} {args}", input=stdin)

    def __str__(self, cmd, *args, **kwargs):
        args = " ".join(build_cli_args_safe(cmd, *args, **kwargs))
        return f"{self.cmd} {args}"


class Genesis:
    def __init__(self, path: Path):
        self.path = path
        self.config = defaultdict(lambda: defaultdict(dict))
        self.data = self.load()

    def load(self):
        "load config from toml file from self.path"
        try:
            with open(self.path, "r") as file:
                self.config.update(json.load(file))
        except FileNotFoundError:
            print(f"Config file {self.path} not found.")

    def save(self):
        "save config to self.path file in toml format"
        with open(self.path, "w") as file:
            json.dump(self.config, file)

    def save_to(self, path):
        "save config to self.path file in toml format"
        with open(path, "w") as file:
            json.dump(self.config, file)

    def deep_update(self, original, new):
        "Recursive update of nested dictionaries"
        for key, value in new.items():
            if isinstance(value, dict) and isinstance(original.get(key), dict):
                self.deep_update(original[key], value)
            else:
                original[key] = value

    def edit(self, new_config):
        "edit config with new_config, keeping unchanged keys"
        self.deep_update(self.config, new_config)
        self.save()


class GalaCLI:
    "the apis to interact with wallet and blockchain"

    def __init__(
        self,
        cmd=DEFAULT_CHAIN_BINARY,
        data_dir=None,
        node_rpc=None,
        chain_id=DEFAULT_TEST_CHAINID,
        keyring_backend="test",
        broadcast_mode="sync",
        output_format="json",  # --output="json"
    ):
        if data_dir:
            self.data_dir = Path(data_dir)
        else:
            temp = tempfile.TemporaryDirectory(delete=True)
            self.data_dir = Path(temp.name)
        self.chain_id = chain_id
        self.keyring_backend = keyring_backend
        self.node_rpc = node_rpc
        self.cmd = cmd
        self.raw = BinaryCommand(cmd)
        self.output = None
        self.output_format = output_format
        self.broadcast_mode = broadcast_mode
        self.error = None
        self.config = None
        self.app_config = None
        self.client_config = None
        self.load_config()

    def load_config(self):
        if Path(self.data_dir / "config/config.toml").exists():
            self.config = GalaToml(self.data_dir / "config/config.toml")
        if Path(self.data_dir / "config/app.toml").exists():
            self.app_config = GalaToml(self.data_dir / "config/app.toml")
        if Path(self.data_dir / "config/client.toml").exists():
            self.client_config = GalaClientConfig(self.data_dir / "config/client.toml")

    def version(self):
        return self.raw("version", home=self.data_dir)

    @property
    def node_rpc_http(self):
        return "http" + self.node_rpc.removeprefix("tcp")

    def status(self):
        return json.loads(self.raw("status", node=self.node_rpc))

    def block_height(self):
        return int(self.status()["sync_info"]["latest_block_height"])

    def block_time(self):
        return isoparse(self.status()["sync_info"]["latest_block_time"])

    def wait_for_block(self, height, timeout=240):
        for _ in range(timeout * 2):
            current_height = self.block_height()
            if current_height >= height:
                return True
            print("current block height", current_height)
            time.sleep(0.5)
        else:
            raise TimeoutError(f"wait for block {height} timeout")

    def rollback(self, hard=False):
        return self.raw("rollback", "--hard" if hard else None, home=self.data_dir)

    ##############################
    #       GENESIS cmds
    ##############################

    def validate_genesis(self):
        return self.raw("validate-genesis", home=self.data_dir)

    def add_genesis_account(self, addr, coins, **kwargs):
        return self.raw(
            "add-genesis-account",
            addr,
            coins,
            home=self.data_dir,
            output="json",
            **kwargs,
        )

    def gentx(self, name, coins, min_self_delegation=1, pubkey=None, **kwargs):
        return self.raw(
            "gentx",
            name,
            coins,
            min_self_delegation=str(min_self_delegation),
            home=self.data_dir,
            chain_id=self.chain_id,
            keyring_backend=self.keyring_backend,
            pubkey=pubkey,
            **kwargs,
        )

    def collect_gentxs(self, gentx_dir):
        return self.raw("collect-gentxs", gentx_dir, home=self.data_dir)

    ##############################
    #     ACCOUNT KEYS utils
    ##############################

    def migrate_keystore(self):
        return self.raw("keys", "migrate", home=self.data_dir)

    def address(self, name, bech="acc"):
        output = self.raw(
            "keys",
            "show",
            name,
            "-a",
            home=self.data_dir,
            keyring_backend=self.keyring_backend,
            bech=bech,
        )
        return output.strip().decode()

    def create_account(self, name, mnemonic=None):
        "create new keypair in node's keyring"
        if mnemonic is None:
            output = self.raw(
                "keys",
                "add",
                name,
                home=self.data_dir,
                algo="eth_secp256k1",
                output="json",
                keyring_backend=self.keyring_backend,
            )
        else:
            output = self.raw(
                "keys",
                "add",
                name,
                "--recover",
                home=self.data_dir,
                algo="eth_secp256k1",
                output="json",
                keyring_backend=self.keyring_backend,
                stdin=mnemonic.encode() + b"\n",
            )
        return json.loads(output)

    def delete_account(self, name):
        "delete wallet account in node's keyring"
        return self.raw(
            "keys",
            "delete",
            name,
            "-y",
            "--force",
            home=self.data_dir,
            output="json",
            keyring_backend=self.keyring_backend,
        )

    ##############################
    #   Tendermint => Cometbft
    ##############################

    def consensus_address(self) -> str:
        "get comet consensus address"
        output = self.raw("comet", "show-address", home=self.data_dir)
        return output.decode().strip()

    def node_id(self) -> str:
        "get comet node id"
        output = self.raw("comet", "show-node-id", home=self.data_dir)
        return output.decode().strip()

    def export(self):
        return self.raw("export", home=self.data_dir)

    def unsaferesetall(self):
        return self.raw("unsafe-reset-all")

    ##############################
    #       FEEMARKET Module
    ##############################

    def query_base_fee(self, **kwargs):
        default_kwargs = {"home": self.data_dir}

        # TODO: is this assumption correct? Having the base fee turned off has caused some test failures
        # because it was returning `null` and not an `int(...)` -> we'll return 0 here.
        params = json.loads(
            self.raw(
                "q",
                "feemarket",
                "params",
                output="json",
                node=self.node_rpc,
                **(default_kwargs | kwargs),
            )
        )
        no_base_fee = params["params"]["no_base_fee"]
        if no_base_fee:
            return 0

        base_fee_out = self.raw(
            "q",
            "feemarket",
            "base-fee",
            output="json",
            node=self.node_rpc,
            **(default_kwargs | kwargs),
        )

        out_dict = json.loads(base_fee_out)
        if not out_dict:
            raise ValueError(f"failed to return base fee: {out_dict}")

        base_fee = out_dict["base_fee"]
        if not base_fee:
            raise ValueError(f"failed to return base fee: {out_dict}")

        return float(base_fee)


class GalaNodeCLI(GalaCLI):
    "Class to control started node of galacticad"

    def __init__(
        self,
        cmd=DEFAULT_CHAIN_BINARY,
        data_dir=None,
        node_rpc=None,
        # node_api=None,
        chain_id=None,
        # node_id=None,
        moniker=DEFAULT_TEST_MONIKER,
        keyring_backend="test",
        node_addr="127.0.0.1",
    ):
        super().__init__(
            cmd=cmd,
            data_dir=data_dir,
            node_rpc=node_rpc,
            chain_id=chain_id,
            keyring_backend=keyring_backend,
        )
        # self.node_id = node_id
        self.moniker = moniker
        self.node_addr = node_addr
        self.account = None
        self.node_rpc = f"tcp://{self.node_addr}:26657"
        self.process = None

    def initial_configure_node(self):
        self.raw("config", "set", "client", "keyring-backend", "test")
        self.load_config()
        self.client_config.edit({"output": "json", "chain-id": self.chain_id})
        self.config.edit({"moniker": self.moniker})
        ## set self network addr
        self.config.apply_addr(self.node_addr)
        self.client_config.apply_addr(self.node_addr)
        self.app_config.apply_addr(self.node_addr)
        ## other initial config
        self.app_config.edit({"api": {"enable": True}})
        self.app_config.edit({"pruning": "nothing"})
        self.app_config.edit({"minimum-gas-prices": f"10{DEFAULT_DENOM}"})
        self.app_config.edit(
            {
                "telemetry": {
                    "service-name": "galacticad",
                    "enabled": True,
                    "prometheus-retention-time": "60",
                    "global-labels": [["chain-id", self.chain_id]],
                }
            }
        )
        self.config.edit(
            {"moniker": self.moniker, "log_format": "json", "log_level": "debug"}
        )
        self.config.edit({"consensus": {"timeout_commit": "1s"}})
        self.config.edit(
            {
                "rpc": {
                    "cors_allowed_origins": [
                        "*",
                    ]
                }
            }
        )

    async def start(self):
        await self.run("start", home=self.data_dir, chain_id=self.chain_id)

    def node_info(self):
        return requests.get(
            f"{self.node_rpc_http}/cosmos/staking/v1beta1/validators/{self.node_id}"
        ).json()

    def init_node(self, moniker=None):
        "### Generate initial config with genesis.json"
        moniker = moniker or self.moniker or DEFAULT_TEST_MONIKER
        return self.raw(
            "init",
            moniker,
            chain_id=self.chain_id,
            home=self.data_dir,
        )

    async def run(self, cmd, *args, **kwargs):
        cmd_args = build_cli_args_safe(cmd, *args, **kwargs)
        with open(self.data_dir / "stdout.log", "a") as out, open(
            self.data_dir / "stderr.log", "a"
        ) as err:
            self.process = await asyncio.create_subprocess_exec(
                self.cmd,
                *cmd_args,
                stdout=out,
                stderr=err,
            )
            return self

    async def get_output(self):
        stdout, stderr = await self.process.communicate()
        return stdout.decode(), stderr.decode(), self.process.returncode

    def is_running(self):
        return self.process and self.process.returncode is None

    async def terminate(self, timeout=30):
        if self.is_running():
            self.process.terminate()
            try:
                await asyncio.wait_for(self.process.wait(), timeout=timeout)
                exit_code = self.process.returncode
                print(f"Instance {self.moniker} exited with code {exit_code}")
                return exit_code
            except asyncio.TimeoutError:
                print(
                    f"Process exceeded timeout of {timeout} seconds. Killing the process."
                )
                self.process.kill()
                await self.process.wait()
                return self.process.returncode

    def set_address_in_configs(self, addr: str):
        for c in (self.client_config, self.app_config, self.config):
            if c:
                c.apply_addr(addr)


class GalaNetwork:
    "### Bunch of GalaNodes with some similar parameters"

    def __init__(self, n_nodes=3, chain_id=DEFAULT_TEST_CHAINID, *args, **kwargs):
        self.chain_id = chain_id
        self.nodes: List[GalaNodeCLI] = []
        self.command_node = GalaNodeCLI(
            data_dir="node00", moniker="node00", chain_id=chain_id
        )
        for n in range(n_nodes):
            name = f"node0{n + 1}"
            self.nodes.append(
                GalaNodeCLI(
                    moniker=name,
                    data_dir=name,
                    node_addr=f"127.0.0.{ 2 + n }",  ## 127.0.0.2 127.0.0.3 ...
                    chain_id=self.chain_id,
                    *args,
                    **kwargs,
                )
            )

    async def initial_configure_network(self):
        for n in self.nodes:
            await n.initial_configure_node()

    async def start(self):
        "### Start every node in network"
        for node in self.nodes:
            await node.start()

    async def stop(self):
        "### Stop every node in network"
        for node in self.nodes:
            await node.terminate()

    async def check_live(self):
        "### Check that node is running"
        ...

    def configure_genesis(self):
        # total_supply = str(GnetAmount("200gnet")) ## will be calculated later
        staking_min_deposit = GnetAmount("100gnet").min_denom_amount_str()
        max_deposit_period = "600s"
        unbonding_time = "30s"
        ## faucet initialzed here because its address neede in genesis minting config
        self.faucet = self.command_node.create_account(
            "faucet", PREDEFINED_KEY_MNEMONIC_FAUCET
        )
        faucet_address = self.faucet["address"]
        inflation_validators_share = "0.99933"
        inflation_faucet_share = "0.00067"

        block_max_gas = "40000000"
        block_max_bytes = "22020096"
        time_iota_ms = "1000"
        voting_period = "60s"
        expedited_voting_period = "30s"
        genesis = Genesis(path=self.command_node.data_dir / "config/genesis.json")
        genesis.edit(
            {
                "consensus": {
                    "params": {
                        "block": {
                            "max_bytes": block_max_bytes,
                            "max_gas": block_max_gas,
                            "time_iota_ms": time_iota_ms,
                        }
                    }
                }
            }
        )
        genesis.edit(
            {"app_state": {"gov": {"voting_params": {"voting_period": voting_period}}}}
        )

        update_genesis = {
            "app_state": {
                "gov": {
                    "deposit_params": {
                        "min_deposit": [
                            {"denom": BASE_DENOM, "amount": staking_min_deposit}
                        ]
                    },
                    "params": {
                        "min_deposit": [
                            {"denom": BASE_DENOM, "amount": staking_min_deposit}
                        ],
                        "max_deposit_period": max_deposit_period,
                        "voting_period": voting_period,
                        "expedited_voting_period": expedited_voting_period,
                    },
                },
                "staking": {
                    "params": {
                        "bond_denom": BASE_DENOM,
                        "unbonding_time": unbonding_time,
                    }
                },
                "crisis": {"constant_fee": {"denom": BASE_DENOM}},
                "mint": {"params": {"mint_denom": BASE_DENOM}},
                "bank": {
                    "denom_metadata": [
                        {
                            "description": "The native staking token of the Galactica Network.",
                            "denom_units": [
                                {
                                    "denom": BASE_DENOM,
                                    "exponent": 0,
                                    "aliases": ["attognet"],
                                },
                                {
                                    "denom": "ugnet",
                                    "exponent": 6,
                                    "aliases": ["micrognet"],
                                },
                                {"denom": DISPLAY_DENOM, "exponent": 18},
                            ],
                            "base": BASE_DENOM,
                            "display": DISPLAY_DENOM,
                            "name": "Galactica Network",
                            "symbol": DISPLAY_DENOM.upper(),
                            "uri": "",
                            "uri_hash": "",
                        }
                    ],
                    "send_enabled": [{"denom": BASE_DENOM, "enabled": True}],
                    # "supply": [{"denom": BASE_DENOM, "amount": total_supply}],
                },
                "evm": {"params": {"evm_denom": BASE_DENOM}},
                "inflation": {
                    "params": {
                        "enable_inflation": True,
                        "mint_denom": BASE_DENOM,
                        "inflation_distribution": {
                            "validators_share": inflation_validators_share,
                            "other_shares": [
                                {
                                    "address": faucet_address,
                                    "name": "faucet",
                                    "share": inflation_faucet_share,
                                }
                            ],
                        },
                    },
                    "inflation_distribution": {
                        "validators_share": inflation_validators_share,
                        "other_shares": [
                            {
                                "address": faucet_address,
                                "name": "faucet",
                                "share": inflation_faucet_share,
                            }
                        ],
                    },
                },
            },
        }
        genesis.edit(update_genesis)

    def combine_seeds(self) -> Dict:
        """
        ### Получение строк которые лягут в seeds или в persistent_peers
        """
        not_combined = {
            node.moniker: node.node_id() + "@" + node.node_addr + ":26656"
            for node in self.nodes
        }
        combined = {
            node.moniker: ",".join(
                [not_combined[n] for n in not_combined if n != node.moniker]
            )
            for node in self.nodes
        }
        return combined

    def configure_network(self):
        """
        ### func for configuring gala network

        - [X] Init first node to get blank genesis.json
        - [X] Edit config files via first node to set common config
        - [X] Configure genesis to needed state
            - [X] edit genesis.json
            - [X] add some gentx
            - [X] collect gentx
            - [X] validate genesis
        - [ ] Get tendermint node-id of each node
            - [X] Put node folder
            - [X] configure node
            - [X] put key into node
            - [X] put genesis.json to node config
            - [ ] get tendermint node-id
        - [ ] Edit individual configs to set some parameters throgh network
            - [ ] Persistent peers
        """

        self.command_node.init_node(
            moniker=self.command_node.moniker,
        )
        self.configure_genesis()

        ## Create treasury account
        treasury_acc = self.command_node.create_account(
            "treasury", mnemonic=PREDEFINED_KEY_MNEMONIC_TREASURY
        )
        if not self.faucet:
            faucet_acc = self.command_node.create_account(
                "faucet", mnemonic=PREDEFINED_KEY_MNEMONIC_FAUCET
            )
        else:
            faucet_acc = self.faucet
        total_supply = GnetAmount("0gnet")
        self.command_node.add_genesis_account(
            treasury_acc["address"], GnetAmount("800gnet")
        )
        self.command_node.add_genesis_account(
            faucet_acc["address"], GnetAmount("100gnet")
        )

        total_supply += GnetAmount("800gnet") + GnetAmount("100gnet")
        self.genesis = Genesis(self.command_node.data_dir / "config/genesis.json")

        ## Create accounts for nodes
        for num, node in enumerate(self.nodes):
            node.init_node()
            self.genesis.save_to(node.data_dir / "config/genesis.json")

            node_genesis_supply = GnetAmount("200gnet")
            node.account = node.create_account(
                node.moniker, mnemonic=PREDEFINED_KEY_MNEMONIC_NODE_KEYS[num]
            )
            node.add_genesis_account(node.account["address"], node_genesis_supply)
            self.command_node.add_genesis_account(
                node.account["address"], node_genesis_supply
            )
            node.gentx(
                node.moniker,
                node_genesis_supply,
                ip=node.node_addr,
                commission_rate="0.02",
                details=f"test.{node.moniker}.details",
            )

            total_supply += node_genesis_supply

        gentx_dir = self.command_node.data_dir / "config/gentx"

        os.makedirs(gentx_dir, exist_ok=True)

        for node in self.nodes:
            node_gentx_dir = node.data_dir / "config/gentx"
            for file_name in os.listdir(node_gentx_dir):
                if file_name.endswith(".json"):
                    file_path = node_gentx_dir / file_name
                    dest_path = gentx_dir / file_name
                    shutil.move(file_path, dest_path)

        self.genesis.load()
        self.genesis.edit(
            {
                "app_state": {
                    "bank": {
                        "supply": [
                            {
                                "denom": BASE_DENOM,
                                "amount": total_supply.min_denom_amount_str(),
                            }
                        ],
                    },
                },
            }
        )
        self.command_node.collect_gentxs(self.command_node.data_dir / "config/gentx")
        self.command_node.validate_genesis()
        self.genesis.load()

        ## configure nodes
        for num, node in enumerate(self.nodes):
            self.genesis.save_to(node.data_dir / "config/genesis.json")
            node.load_config()
            node.initial_configure_node()
            node.config.apply_addr(node.node_addr)

        combine_seeds = self.combine_seeds()
        for node in self.nodes:
            node.config.edit({"p2p": {"persistent_peers": combine_seeds[node.moniker]}})


async def main():
    chain_id = DEFAULT_TEST_CHAINID
    g_network = GalaNetwork(3, chain_id=chain_id)
    g_network.configure_network()

    print("start node...")
    await g_network.start()
    await asyncio.sleep(5)

    g_network.nodes[0].wait_for_block(5)
    await g_network.stop()

    # await asyncio.sleep(1)
    # await gn1.start()
    # await asyncio.sleep(1)
    # print(gn1.is_running())
    # gn1.wait_for_block(5, timeout=30)
    # # await gn1.terminate()
    # # await asyncio.sleep(10)
    # if await gn1.terminate() == 0:
    #     print("Success")


if __name__ == "__main__":
    asyncio.run(main())
