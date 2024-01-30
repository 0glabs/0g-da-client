#!/usr/bin/env python3
import argparse
import os
import subprocess
import sys

sys.path.append("../zerog_storage_kv/tests")

from concurrent.futures import ProcessPoolExecutor
from test_all import run_single_test

PORT_MIN = 11000
PORT_MAX = 65535
PORT_RANGE = 600

__file_path__ = os.path.dirname(os.path.realpath(__file__))


def run():
    parser = argparse.ArgumentParser(usage="%(prog)s [options]")
    parser.add_argument(
        "--max-workers",
        dest="max_workers",
        default=5,
        type=int,
    )
    parser.add_argument(
        "--port-max",
        dest="port_max",
        default=PORT_MAX,
        type=int,
    )
    parser.add_argument(
        "--port-min",
        dest="port_min",
        default=PORT_MIN,
        type=int,
    )

    options = parser.parse_args()

    TEST_SCRIPTS = []

    test_dir = os.path.dirname(os.path.realpath(__file__))
    test_subdirs = [
        "",  # include test_dir itself
    ]

    slow_tests = {}

    for subdir in test_subdirs:
        subdir_path = os.path.join(test_dir, subdir)
        for file in os.listdir(subdir_path):
            if file.endswith("_test.py"):
                rel_path = os.path.join(subdir, file)
                if rel_path not in slow_tests:
                    TEST_SCRIPTS.append(rel_path)

    executor = ProcessPoolExecutor(max_workers=options.max_workers)
    test_results = []

    py = "python"
    if hasattr(sys, "getwindowsversion"):
        py = "python"

    i = 0
    # Start slow tests first to avoid waiting for long-tail jobs
    for script in slow_tests:
        f = executor.submit(
            run_single_test, py, script, test_dir, i, options.port_min, options.port_max
        )
        test_results.append((script, f))
        i += 1
    for script in TEST_SCRIPTS:
        f = executor.submit(
            run_single_test, py, script, test_dir, i, options.port_min, options.port_max
        )
        test_results.append((script, f))
        i += 1

    failed = set()
    for script, f in test_results:
        try:
            f.result()
        except subprocess.CalledProcessError as err:
            print("CalledProcessError " + repr(err))
            failed.add(script)

    if len(failed) > 0:
        print("The following test fails: ")
        for c in failed:
            print(c)
        sys.exit(1)


if __name__ == "__main__":
    run()
