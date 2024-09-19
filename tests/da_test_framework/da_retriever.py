import os
import sys
import time
import grpc
import retriever_pb2 as pb2
import retriever_pb2_grpc as pb2_grpc
from da_test_framework.da_node_type import DANodeType

sys.path.append("../0g-storage-kv/tests")

from test_framework.blockchain_node import TestNode


__file_path__ = os.path.dirname(os.path.realpath(__file__))


class DARetriever(TestNode):
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

        data_dir = os.path.join(root_dir, "da_retriever")
        self.grpc_url = "0.0.0.0:32011"
        super().__init__(
            DANodeType.DA_RETRIEVER,
            14,
            data_dir,
            None,
            binary,
            local_conf,
            log,
            None,
        )
        self.args = [
            binary,
            "--retriever.hostname", "localhost",
            "--retriever.grpc-port", "32011",
            "--retriever.storage.node-url", f'http://{local_conf["node_rpc_endpoint"]}',
            "--retriever.storage.kv-url", f'http://{local_conf["kv_rpc_endpoint"]}',
            "--retriever.storage.kv-stream-id", local_conf['stream_id'],
            "--retriever.storage.flow-contract", local_conf['log_contract_address'],
            "--retriever.log.level-std", "trace",
            "--kzg.srs-order", "300000",
        ]
    
    def wait_for_rpc_connection(self):
        # TODO: health check of service availability
        time.sleep(3)
        self.channel = grpc.insecure_channel(self.grpc_url)
        # bind the client and the server
        self.stub = pb2_grpc.RetrieverStub(self.channel)
       
    def start(self):
        self.log.info("Start DA retriever")
        super().start()

    def stop(self):
        self.log.info("Stop DA retriever")
        try:
            super().stop(kill=True, wait=False)
        except AssertionError as e:
            err = repr(e)
            if "no RPC connection" in err:
                self.log.debug(f"Stop DA retriever: no RPC connection")
            else:
                raise e

    def retrieve_blob(self, info):
        message = pb2.BlobRequest(batch_header_hash=info.blob_verification_proof.batch_metadata.batch_header_hash, blob_index=info.blob_verification_proof.blob_index)
        return self.stub.RetrieveBlob(message)
