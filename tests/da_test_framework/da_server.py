import os
import sys
import time

sys.path.append("../0g-storage-kv/tests")

from test_framework.blockchain_node import TestNode, FailedToStartError
from da_test_framework.da_node_type import DANodeType
from utility.simple_rpc_proxy import SimpleRpcProxy

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
        print(updated_config)
        print(local_conf)
        local_conf.update(updated_config)

        data_dir = os.path.join(root_dir, "da_server")
        rpc_url = "http://0.0.0.0:51001"
        super().__init__(
            DANodeType.DA_SERVER,
            13,
            data_dir,
            rpc_url,
            binary,
            local_conf,
            log,
            10,
        )
        self.args = [binary, "--disperser-server.grpc-port", "51001",
                     "--disperser-server.s3-bucket-name", "test-zgda-blobstore",
                     "--disperser-server.dynamodb-table-name", "test-BlobMetadata",
                     "--disperser-server.aws.region", "us-east-1",
                     "--disperser-server.aws.access-key-id", "localstack",
                     "--disperser-server.aws.secret-access-key", "localstack",
                     "--disperser-server.aws.endpoint-url", "http://0.0.0.0:4566"]
    
    def wait_for_rpc_connection(self):
        self._wait_for_rpc_connection(lambda rpc: True)

    def start(self):
        self.log.info("Start DA server")
        super().start()

    def stop(self):
        self.log.info("Stop DA server")
        super().stop(kill=True, wait=False)
        
    def disperse_blob(self, request):
        return self.rpc.DisperseBlob(request)

    def retrieve_blob(self, request):
        return self.rpc.RetrieveBlob(request)

    def get_blob_status(self, request_id):
        return self.rpc.GetBlobStatus({'request_id': request_id})
