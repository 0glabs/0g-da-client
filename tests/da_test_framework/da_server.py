import os
import sys
import time

sys.path.append("../../zerog_storage_kv/tests")

from test_framework.blockchain_node import TestNode
from da_test_framework.da_node_type import DANodeType

__file_path__ = os.path.dirname(os.path.realpath(__file__))


class DAServer(TestNode):
    def __init__(
            self,
            root_dir,
            binary,
            updated_config,
            log,
    ):
        local_conf = dict(log_config_file="log_config")

        local_conf.update(updated_config)
        data_dir = os.path.join(root_dir, "da_server")
        # rpc_url = "http://" + local_conf["rpc_listen_address"]
        super().__init__(
            DANodeType.DA_SERVER,
            0,
            data_dir,
            None,
            binary,
            local_conf,
            log,
            None,
        )
        self.args = [binary, "--disperser-server.grpc-port", "51001",
                     "--disperser-server.s3-bucket-name", "test-zgda-blobstore",
                     "--disperser-server.dynamodb-table-name", "test-BlobMetadata",
                     "--disperser-server.aws.region", "us-east-1",
                     "--disperser-server.aws.access-key-id", "localstack",
                     "--disperser-server.aws.secret-access-key", "localstack",
                     "--disperser-server.aws.endpoint-url", "http://0.0.0.0:4566"]

    def start(self):
        self.log.info("Start DA server")
        super().start()

    def wait_for_rpc_connection(self):
        time.sleep(1)

    def stop(self):
        self.log.info("Stop DA server")
        super().stop(kill=True, wait=False)
