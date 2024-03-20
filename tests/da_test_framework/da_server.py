import os
import sys
import time
import grpc
import disperser_pb2 as pb2
import disperser_pb2_grpc as pb2_grpc
from da_test_framework.da_node_type import DANodeType

sys.path.append("../0g-storage-kv/tests")

from test_framework.blockchain_node import TestNode, FailedToStartError
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
        self.grpc_url = "0.0.0.0:51001"
        super().__init__(
            DANodeType.DA_SERVER,
            13,
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
    
    def wait_for_rpc_connection(self):
        poll = self.process.poll()
        if poll is None:
            time.sleep(3)
        self.log.info(f'{self.grpc_url}')
        self.channel = grpc.insecure_channel(self.grpc_url)
        # bind the client and the server
        self.stub = pb2_grpc.DisperserStub(self.channel)
       
    def start(self):
        self.log.info("Start DA server")
        super().start()

    def stop(self):
        self.log.info("Stop DA server")
        super().stop(kill=True, wait=False)
        
    def disperse_blob(self, data):
        message = pb2.DisperseBlobRequest(data=data, security_params=[pb2.SecurityParams(quorum_id=0, adversary_threshold=25, quorum_threshold=50)], target_chunk_num=16)
        return self.stub.DisperseBlob(message)

    def retrieve_blob(self, info):
        message = pb2.RetrieveBlobRequest(batch_header_hash=info.blob_verification_proof.batch_metadata.batch_header_hash, blob_index=info.blob_verification_proof.blob_index)
        self.log.info(f'retrieve blob {message}')
        return self.stub.RetrieveBlob(message)

    def get_blob_status(self, request_id):
        message = pb2.BlobStatusRequest(request_id=request_id)
        self.log.info(f'get blob status {message}')
        return self.stub.GetBlobStatus(message)
