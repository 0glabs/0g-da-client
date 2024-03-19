#!/usr/bin/env python3
import random
import sys
sys.path.append("../0g_storage_kv/tests")
import time

from da_test_framework.da_test_framework import DATestFramework
from utility.utils import assert_equal


class DAPutGetTest(DATestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 2

    def run_test(self):
        # setup kv node, watch stream with id [0,100)
        request = self.create_request()
        client = self.da_services[-1]
        reply = client.disperse_blob(request)
        while client.get_blob_status(reply.RequestId).status != 2:
            time.sleep(1)
        # retrieve the blob
        reply = client.retrieve_blob(request)
        assert_equal(reply.data, request.data)
        
    def create_request():
        data = [random.randint(0, 255) for _ in range(100)]
        return {
            'data': data.to_vec(),
            'security_params': [
                {
                    'quorum_id': 0,
                    'adversary_threshold': 25,
                    'quorum_threshold': 50,
                }
            ],
            'target_chunk_num': 32,
        }


if __name__ == "__main__":
    DAPutGetTest().main()
