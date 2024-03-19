import os
import sys
import time

sys.path.append("../../0g-storage-kv/tests")

from test_framework.blockchain_node import TestNode
from da_test_framework.da_node_type import DANodeType

__file_path__ = os.path.dirname(os.path.realpath(__file__))


class DAEncoder(TestNode):
    def __init__(
            self,
            root_dir,
            binary,
            updated_config,
            log,
    ):
        local_conf = dict(log_config_file="log_config")

        local_conf.update(updated_config)
        data_dir = os.path.join(root_dir, "da_encoder")
        # rpc_url = "http://" + local_conf["rpc_listen_address"]
        super().__init__(
            DANodeType.DA_ENCODER,
            0,
            data_dir,
            None,
            binary,
            local_conf,
            log,
            None,
        )
        self.args = [binary, "--disperser-encoder.grpc-port", "34000",
                     "--disperser-encoder.metrics-http-port", "9109",
                     "--kzg.g1-path", f"{__file_path__}/../../inabox/resources/kzg/g1.point.300000",
                     "--kzg.g2-path", f"{__file_path__}/../../inabox/resources/kzg/g2.point.300000",
                     "--kzg.cache-path", f"{__file_path__}/../../inabox/resources/kzg/SRSTables",
                     "--kzg.srs-order", "300000",
                     "--kzg.num-workers", "12",
                     "--disperser-encoder.log.level-std", "trace",
                     "--disperser-encoder.log.level-file", "trace"]

    def start(self):
        self.log.info("Start DA encoder")
        super().start()

    def wait_for_rpc_connection(self):
        time.sleep(1)

    def stop(self):
        self.log.info("Stop DA encoder")
        # The encoder will check return_code via rpc when error log exists
        # that is written when the encoder starts normally.
        # The exception handling can be removed when rpc is added or the error
        # is not written when the encoder starts normally.
        try:
            super().stop(kill=True, wait=False)
        except AssertionError as e:
            err = repr(e)
            if "no RPC connection" in err:
                self.log.debug(f"Stop DA encoder: no RPC connection")
            else:
                raise e
