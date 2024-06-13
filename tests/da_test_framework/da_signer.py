import os
import sys
import time


from da_test_framework.test_node import TestNode
from da_test_framework.da_node_type import DANodeType

from config.node_config import PRIV_KEY


__file_path__ = os.path.dirname(os.path.realpath(__file__))


class DASigner(TestNode):
    def __init__(
        self,
        root_dir,
        binary,
        updated_config,
        log,
    ):
        local_conf = dict(
            log_level = "info",
            data_path = "./db/",
            da_entrance_address = "0x64fcfde2350E08E7BaDc18771a7674FAb5E137a2",
            start_block_number = 0,
            signer_private_key = "1",
            validator_private_key = PRIV_KEY,
        )

        local_conf.update(updated_config)
        data_dir = os.path.join(root_dir, "da_signer")
        super().__init__(
            DANodeType.DA_SIGNER,
            11,
            data_dir,
            None,
            binary,
            local_conf,
            log,
            None,
        )

    def start(self):
        self.log.info("Start DA signer")
        super().start()

    def wait_for_rpc_connection(self):
        time.sleep(3)

    def stop(self):
        self.log.info("Stop DA signer")

        try:
            super().stop(kill=True, wait=False)
        except AssertionError as e:
            err = repr(e)
            if "no RPC connection" in err:
                self.log.debug(f"Stop DA signer: no RPC connection")
            else:
                raise e
