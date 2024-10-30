import subprocess
from dateutil.parser import isoparse
import json
import time
import socket
import os
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
DEFAULT_CHAIN_BINARY = "../../build/galacticad"

DEFAULT_TEST_MONIKER = "test-node01"
DEFAULT_TEST_CHAINID = "test_41239-41239"

PREDEFINED_KEY_MNEMONIC = "gesture inject test cycle original hollow east ridge hen combine junk child bacon zero hope comfort vacuum milk pitch cage oppose unhappy lunar seat"


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
    (stdout, _) = proc.communicate(input=input)
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
        "load config from toml file from self.path"
        try:
            with open(self.path, "r") as file:
                self.config = toml.load(file)
        except FileNotFoundError:
            print(f"Config file {self.path} not found.")

    def save(self):
        "save config to self.path file in toml format"
        with open(self.path, "w") as file:
            toml.dump(self.config, file)

    def edit(self, new_config):
        "edit config with new_config"
        for key, value in new_config.items():
            if isinstance(value, dict):
                self.config[key].update(value)
            else:
                self.config[key] = value
        self.save()

    def diff(self, other):
        "return a dictionary containing the differences between self and other"
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
        "helper method to compare two dictionaries"
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
        "replace the host in all address values in config with the host from new_addr"
        new_host = urllib.parse.urlparse(new_addr).hostname or new_addr
        self._apply_addr_to_dict(self.config, new_host)
        self.save()

    def _apply_addr_to_dict(self, inner_dict, new_host):
        "helper method to apply new_host to nested dictionaries recursively"
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
        # self.raw("config", home=self.data_dir, chain_id=self.chain_id)
        # self.config = GalaToml(self.data_dir / "config/config.toml")
        # self.app_config = GalaToml(self.data_dir / "config/app.toml")
        # self.client_config = GalaClientConfig(self.data_dir / "config/client.toml")

    @property
    def node_rpc_http(self):
        return "http" + self.node_rpc.removeprefix("tcp")

    def status(self):
        return json.loads(self.raw("status", node=self.node_rpc))

    def block_height(self):
        return int(self.status()["sync_info"]["latest_block_height"])

    def block_time(self):
        return isoparse(self.status()["sync_info"]["latest_block_time"])

    def rollback(self, hard=False):
        return self.raw("rollback", "--hard" if hard else None, home=self.data_dir)

    def version(self):
        return self.raw("version", home=self.data_dir)

    # ==========================
    #       GENESIS cmds
    # ==========================

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

    def gentx(self, name, coins, min_self_delegation=1, pubkey=None):
        return self.raw(
            "gentx",
            name,
            coins,
            min_self_delegation=str(min_self_delegation),
            home=self.data_dir,
            chain_id=self.chain_id,
            keyring_backend=self.keyring_backend,
            pubkey=pubkey,
        )

    def collect_gentxs(self, gentx_dir):
        return self.raw("collect-gentxs", gentx_dir, home=self.data_dir)

    # ==========================
    #     ACCOUNT KEYS utils
    # ==========================

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


class GalaNodeCLI(GalaCLI):
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
        self.process = None

    def start(self): ...

    def node_info(self):
        return requests.get(
            f"{self.node_rpc_http}/cosmos/staking/v1beta1/validators/{self.node_id}"
        ).json()

    def init(self, moniker=None):
        "generate initial config with genesis.json"
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

    async def terminate(self):
        if self.is_running():
            self.process.terminate()
            await self.process.wait()


async def main():
    chain_id = DEFAULT_TEST_CHAINID
    moniker = "test-node01"
    gn1 = GalaNodeCLI(
        data_dir="node01",
        chain_id=chain_id,
        moniker=moniker,
        node_rpc="tcp://127.0.0.2:26657",
    )
    gn1.init()

    gn1.client_config = GalaClientConfig(gn1.data_dir / "./config/client.toml")
    gn1.client_config.edit(
        {
            "chain-id": chain_id,
            "keyring-backend": "test",
            "output": "json",
        }
    )

    gn1.app_config = GalaToml(gn1.data_dir / "./config/app.toml")
    gn1.app_config.apply_addr("127.0.0.2")
    gn1.app_config.edit({"api": {"enable": True}})
    gn1.app_config.edit({"pruning": "nothing"})
    gn1.app_config.edit({"minimum-gas-prices": f"10{DEFAULT_DENOM}"})
    gn1.app_config.edit(
        {
            "telemetry": {
                "service-name": "galacticad",
                "enabled": True,
                "prometheus-retention-time": "60",
                "global-labels": [["chain-id", DEFAULT_TEST_CHAINID]],
            }
        }
    )

    gn1.config = GalaToml(gn1.data_dir / "./config/config.toml")
    gn1.config.edit({"moniker": "test-node01", "log_format": "json"})
    gn1.config.apply_addr("127.0.0.2")
    gn1.config.edit(
        {
            "rpc": {
                "cors_allowed_origins": [
                    "*",
                ]
            }
        }
    )
    ## just check diff between andromeda
    # andr_app_config = GalaToml(
    #     "/Users/booger/occamfi/galactica/ansible-node/gala-andr-root/galaandr_41238-41238/config/app.toml"
    # )
    # d = andr_app_config.diff(app_config)
    print("config created")
    account = gn1.create_account("local-test", PREDEFINED_KEY_MNEMONIC)
    pprint(account)
    pprint(gn1.address("local-test"))

    print("configure genesis")
    await asyncio.sleep(1)
    await gn1.run("start", home=gn1.data_dir, chain_id=gn1.chain_id)
    await asyncio.sleep(1)
    print(gn1.is_running())

    await gn1.terminate()


if __name__ == "__main__":
    asyncio.run(main())
