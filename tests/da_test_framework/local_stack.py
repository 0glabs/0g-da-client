import os
import sys
import time

sys.path.append("../0g_storage_kv/tests")

from test_framework.blockchain_node import TestNode
from da_test_framework.da_node_type import DANodeType


class LocalStack(TestNode):
    def __init__(
            self,
            root_dir,
            binary,
            updated_config,
            log,
    ):
        local_conf = dict(log_config_file="log_config")

        local_conf.update(updated_config)
        data_dir = os.path.join(root_dir, "localstack")
        rpc_url = "http://0.0.0.0:4566"
        super().__init__(
            DANodeType.DA_LOCAL_STACK,
            10,
            data_dir,
            rpc_url,
            binary,
            local_conf,
            log,
            10,
        )
        self.args = [binary, "--localstack-port", "4566", "--deploy-resources", "true", "localstack"]

    def start(self):
        self.log.info("Start localstack")
        super().start()

    def wait_for_rpc_connection(self):
        self._wait_for_rpc_connection(lambda rpc: True)

    def stop(self):
        self.log.info("Stop localstack")
        os.system("docker stop localstack-test")
