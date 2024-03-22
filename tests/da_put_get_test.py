#!/usr/bin/env python3
import random
import sys
from random import randbytes
import time

from disperser_pb2 import BlobStatus

sys.path.append("../0g-storage-kv/tests")
from da_test_framework.da_test_framework import DATestFramework
from utility.utils import assert_equal


class DAPutGetTest(DATestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 1

    def run_test(self):
        client = self.da_services[-1]
        
        data = randbytes(507904)
        disperse_response = client.disperse_blob(data)
        
        self.log.info(disperse_response)
        request_id = disperse_response.request_id
        reply = client.get_blob_status(request_id)
        count = 0
        while reply.status != BlobStatus.CONFIRMED and count <= 5:
            reply = client.get_blob_status(request_id)
            self.log.info(f'blob status {reply.status}')
            count += 1
            time.sleep(10)
        
        info = reply.info
        self.log.info(f'reply info {info}')
        # retrieve the blob
        reply = client.retrieve_blob(info)
        self.log.info(f'reply data {reply.data}')
        assert_equal(reply.data[:len(data)], data)


if __name__ == "__main__":
    DAPutGetTest(
        blockchain_node_configs=dict([(0, dict(mode="dev", dev_block_interval_ms=50))])
    ).main()
