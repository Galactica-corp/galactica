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


def test_010_add_account(gala_network_fixture: GalaNetwork):
    """### Testing"""
