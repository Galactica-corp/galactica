import json
import time
import pytest
from .galacli import (
    DEFAULT_TEST_CHAINID,
    GalaNetwork,
    PREDEFINED_KEY_MNEMONIC_TEST_ACC,
    GnetAmount,
    wait_for_new_blocks,
)
import time
import pytest
from .galacli import DEFAULT_TEST_CHAINID, GalaNetwork


@pytest.fixture(scope="module")
def gala_network_fixture():
    chain_id = DEFAULT_TEST_CHAINID
    g_network = GalaNetwork(3, chain_id=chain_id)

    try:
        yield g_network
    finally:
        g_network.stop_network()
        g_network.clean()


def test_001_gala_network_configure(gala_network_fixture: GalaNetwork):
    gala_network_fixture.configure_network()
    time.sleep(1)
    is_configured = all(node.config for node in gala_network_fixture.nodes)
    assert is_configured


def test_002_gala_network_is_starting(gala_network_fixture: GalaNetwork):
    gala_network_fixture.start_network()
    time.sleep(1)
    is_starting = all(node.is_running() for node in gala_network_fixture.nodes)
    assert is_starting


def test_003_gala_network_is_running(gala_network_fixture: GalaNetwork):
    time.sleep(5)
    is_starting = all(node.status() for node in gala_network_fixture.nodes)
    assert is_starting


def test_004_gala_network_is_growing(gala_network_fixture: GalaNetwork):
    height_before = gala_network_fixture.nodes[0].block_height()
    time.sleep(5)
    height_after = gala_network_fixture.nodes[0].block_height()
    assert height_after > height_before


def test_005_account_add(gala_network_fixture: GalaNetwork):
    node = gala_network_fixture.nodes[2]
    res = node.create_account(
        mnemonic=PREDEFINED_KEY_MNEMONIC_TEST_ACC, name="test_acc"
    )
    assert res["address"] == "gala104ztrup6nxqrlnm4tqh4pxzlzxkw2h4daksgsg"


def test_006_account_list(gala_network_fixture: GalaNetwork):
    node = gala_network_fixture.nodes[2]
    res_json = node.raw(
        "keys", "list", home=node.data_dir, keyring_backend="test", output="json"
    )
    res = json.loads(res_json)
    assert len(res) > 1
    assert res[1]["name"] == "test_acc"


def test_007_transaction_gnet(gala_network_fixture: GalaNetwork):
    cli = gala_network_fixture.command_node
    cli.node_rpc = "tcp://127.0.0.2:26657"
    transfer = cli.transfer(
        "treasury", "gala104ztrup6nxqrlnm4tqh4pxzlzxkw2h4daksgsg", GnetAmount("10gnet")
    )
    tx = transfer["txhash"]
    wait_for_new_blocks(cli, 2)
    tx_res = cli.tx(tx)
    assert tx_res["code"] == 0
    transfer_event = [
        event for event in tx_res["events"] if event["type"] == "transfer"
    ][0]
    assert transfer_event == {
        "type": "transfer",
        "attributes": [
            {
                "key": "recipient",
                "value": "gala17xpfvakm2amg962yls6f84z3kell8c5lq5t8fz",
                "index": True,
            },
            {
                "key": "sender",
                "value": "gala10jmp6sgh4cc6zt3e8gw05wavvejgr5pwqlc3wn",
                "index": True,
            },
            {"key": "amount", "value": "1000000000000000000agnet", "index": True},
        ],
    }


def test_008_bank(gala_network_fixture: GalaNetwork):
    cli = gala_network_fixture.command_node
    cli.node_rpc = "tcp://127.0.0.2:26657"
    balance = cli.balance("gala104ztrup6nxqrlnm4tqh4pxzlzxkw2h4daksgsg")
    assert balance != 0


def test_100_account_del(gala_network_fixture: GalaNetwork):
    node = gala_network_fixture.nodes[2]
    delete_msg = node.delete_account(name="test_acc")
    assert delete_msg == b"Key deleted forever (uh oh!)\n"
    keys_json = node.raw(
        "keys", "list", home=node.data_dir, keyring_backend="test", output="json"
    )
    keys = json.loads(keys_json)
    assert len(keys) == 1
