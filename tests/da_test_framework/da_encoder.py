import os
import sys
import time

sys.path.append("../0g-storage-kv/tests")

from da_test_framework.test_node import TestNode
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
        local_conf = dict(
            log_level = "debug",
        )
        local_conf.update(updated_config)
        data_dir = os.path.join(root_dir, "da_encoder")
        super().__init__(
            DANodeType.DA_ENCODER,
            11,
            data_dir,
            None,
            binary,
            local_conf,
            log,
            None,
        )

    def start(self):
        self.log.info("Start DA encoder")
        super().start()

    def wait_for_rpc_connection(self):
        time.sleep(3)

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
