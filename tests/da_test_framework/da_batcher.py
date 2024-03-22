import os
import sys
import time

sys.path.append("../0g-storage-kv/tests")

from test_framework.blockchain_node import TestNode
from utility.utils import blockchain_rpc_port
from config.node_config import GENESIS_PRIV_KEY
from da_test_framework.da_node_type import DANodeType

__file_path__ = os.path.dirname(os.path.realpath(__file__))


class DABatcher(TestNode):
    def __init__(
            self,
            root_dir,
            binary,
            updated_config,
            log,
    ):
        local_conf = {
            "log_config_file": "log_config",
            "blockchain_rpc_endpoint": f"http://127.0.0.1:{blockchain_rpc_port(0)}",
        }

        local_conf.update(updated_config)
        data_dir = os.path.join(root_dir, "da_batcher")
        super().__init__(
            DANodeType.DA_BATCHER,
            12,
            data_dir,
            None,
            binary,
            local_conf,
            log,
            None,
        )
        self.args = [binary, "--batcher.pull-interval", "10s",
                     "--chain.rpc", local_conf['blockchain_rpc_endpoint'],
                     "--chain.private-key", GENESIS_PRIV_KEY,
                     "--batcher.finalizer-interval", "20s",
                     "--batcher.aws.region", "us-east-1",
                     "--batcher.aws.access-key-id", "localstack",
                     "--batcher.aws.secret-access-key", "localstack",
                     "--batcher.aws.endpoint-url", "http://0.0.0.0:4566",
                     "--batcher.s3-bucket-name", "test-zgda-blobstore",
                     "--batcher.dynamodb-table-name", "test-BlobMetadata",
                     "--encoder-socket", "0.0.0.0:34000",
                     "--batcher.batch-size-limit", "10000",
                     "--batcher.srs-order", "300000",
                     "--encoding-timeout", "10s",
                     "--chain-read-timeout", "12s",
                     "--chain-write-timeout", "13s",
                     "--batcher.storage.node-url", f'http://{local_conf["node_rpc_endpoint"]}',
                     "--batcher.storage.kv-url", f'http://{local_conf["kv_rpc_endpoint"]}',
                     "--batcher.storage.kv-stream-id", local_conf['stream_id'],
                     "--batcher.storage.flow-contract", local_conf['log_contract_address']]

    def start(self):
        self.log.info("Start DA batcher")
        super().start()

    def wait_for_rpc_connection(self):
        time.sleep(15)

    def stop(self):
        self.log.info("Stop DA batcher")
        try:
            super().stop(kill=True, wait=False)
        except AssertionError as e:
            err = repr(e)
            # The batcher will check return_code via rpc when error log exists
            # that is written when the batcher starts normally.
            # The exception handling can be removed when rpc is added or the error
            # is not written when the batcher starts normally.
            if "no RPC connection" in err:
                self.log.debug(f"Stop DA encoder: no RPC connection")
            else:
                raise e
