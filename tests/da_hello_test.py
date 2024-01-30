#!/usr/bin/env python3
import sys
sys.path.append("../zerog_storage_kv/tests")
import time

from da_test_framework.da_test_framework import DATestFramework


class DAHelloTest(DATestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 2

    def run_test(self):
        self.log.debug("===============================================")
        self.log.debug("root dir:" + self.root_dir)
        self.log.debug("blockchain binary:" + self.blockchain_binary)
        self.log.debug("zgs binary:" + self.zgs_binary)
        self.log.debug("cli binary:" + self.cli_binary)
        self.log.debug("kv binary:" + self.kv_binary)
        self.log.debug("contract path:" + self.contract_path)
        self.log.debug("token contract path:" + self.token_contract_path)
        self.log.debug("mine contract path:" + self.mine_contract_path)
        self.log.debug("localstack binary:" + self.localstack_binary)
        self.log.debug("da encoder binary:" + self.da_encoder_binary)
        self.log.debug("da batcher binary:" + self.da_batcher_binary)
        self.log.debug("da server binary:" + self.da_server_binary)
        self.log.debug("===============================================")

        self.log.info("hello")
        time.sleep(3)


if __name__ == "__main__":
    DAHelloTest().main()
