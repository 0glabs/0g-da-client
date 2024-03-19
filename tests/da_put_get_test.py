#!/usr/bin/env python3
import random
import sys
sys.path.append("../0g-storage-kv/tests")
import time

from da_test_framework.da_test_framework import DATestFramework
from utility.utils import assert_equal


class DAPutGetTest(DATestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 1

    def run_test(self):
        # setup kv node, watch stream with id [0,100)
        request = self.create_disperse_request()
        client = self.da_services[-1]
        reply = client.disperse_blob(request)
        request_id = reply.RequestId
        while reply.status != 2:
            time.sleep(10)
            reply = client.get_blob_status(request_id)
        
        info = reply.Info
        # retrieve the blob
        retrieve_request = {
            'BatchHeaderHash':      info.BlobVerificationProof.BatchMetadata.BatchHeaderHash,
            'BlobIndex':            info.BlobVerificationProof.BlobIndex,
            'ReferenceBlockNumber': 0,
            'QuorumId':             0,
        }
        reply = client.retrieve_blob(retrieve_request)
        assert_equal(reply.data[:len(request.data)], request.data)
        
    def create_disperse_request(self):
        data = [random.randint(0, 255) for _ in range(1000)]
        return {
            'data': data,
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
