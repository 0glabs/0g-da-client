import json
import os
import subprocess
import tempfile
import time
import rlp
import shutil

from eth_utils import decode_hex, keccak
from web3 import Web3, HTTPProvider
from enum import Enum, unique
from config.node_config import     PRIV_KEY

from utility.simple_rpc_proxy import SimpleRpcProxy
from utility.utils import (
    initialize_config,
    wait_until,
)
from da_test_framework.contracts import load_contract_metadata

from da_test_framework.test_node import TestNode

@unique
class NodeType(Enum):
    BlockChain = 0
    Zgs = 1
    KV = 2


class FailedToStartError(Exception):
    """Raised when a node fails to start correctly."""


class BlockchainNode(TestNode):
    def __init__(
        self,
        data_dir,
        rpc_url,
        binary,
        local_conf,
        contract_path,
        log,
        rpc_timeout=10,
    ):
        self.contract_path = contract_path

        super().__init__(
            NodeType.BlockChain,
            0,
            data_dir,
            rpc_url,
            binary,
            local_conf,
            log,
            rpc_timeout,
        )

        my_env = os.environ.copy()
        idx = rpc_url.find("://")
        url = rpc_url[idx+3:] if idx != -1 else rpc_url
            
        self.args = [
            binary, "start",
            "--home", os.path.join(my_env["HOME"], ".0gchain"),
            "--json-rpc.address", url,
        ]

    def wait_for_rpc_connection(self):
        time.sleep(10)
        self._wait_for_rpc_connection(lambda rpc: rpc.eth_syncing() is False)

    def wait_for_start_mining(self):
        self._wait_for_rpc_connection(lambda rpc: int(rpc.eth_blockNumber(), 16) > 0)

    def wait_for_transaction_receipt(self, w3, tx_hash, timeout=120, parent_hash=None):
        return w3.eth.wait_for_transaction_receipt(tx_hash, timeout)

    def setup_contract(self):
        origin_path = os.getcwd()
        os.chdir(self.contract_path)

        p = os.path.join(self.contract_path, 'hardhat.config.ts')
        clone_command = "git checkout -- %s" % p
        os.system(clone_command)
        with open(p, 'r') as file:
            file_content = file.read()

        # Replace the string
        modified_content = file_content.replace('http://0.0.0.0:8545', self.rpc_url)

        # Write the modified content back to the file
        with open(p, 'w') as file:
            file.write(modified_content)

        p = os.path.join(self.contract_path, "deployments")
        if os.path.exists(p):
            shutil.rmtree(p)

        p = os.path.join(self.contract_path, ".env")
        if os.path.exists(p):
            os.remove(p)
            
        try:
            with open(p, 'w') as file:
                file.write(f"DEPLOYER_KEY = \"{PRIV_KEY}\"")
        except IOError as e:
            raise e
    
        os.system("yarn")
        os.system("yarn build")
        os.system("yarn deploy zg")

        os.chdir(origin_path)

    def wait_for_transaction(self, tx_hash):
        w3 = Web3(HTTPProvider(self.rpc_url))
        w3.eth.wait_for_transaction_receipt(tx_hash)

    def start(self):
        super().start(False)
