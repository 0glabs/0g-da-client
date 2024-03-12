#!/usr/bin/env python3
import os
import sys

sys.path.append("../0g_storage_kv/0g-storage-node/tests")

from utility.run_all import run_all

if __name__ == "__main__":
    run_all(test_dir=os.path.dirname(__file__))
