import os
import sys
import shutil
import stat
import argparse
import logging
import time
import random
from enum import Enum
import pdb
import tempfile
import subprocess
import traceback

from utility.utils import PortMin, is_windows_platform, blockchain_rpc_port
from da_test_framework.da_encoder import DAEncoder
from da_test_framework.da_server import DAServer
from da_test_framework.blockchain_node import BlockchainNode
from da_test_framework.da_signer import DASigner

__file_path__ = os.path.dirname(os.path.realpath(__file__))
binary_ext = ".exe" if is_windows_platform() else ""

amt_params = "/data/share/amt-params"

class TestStatus(Enum):
    PASSED = 1
    FAILED = 2

TEST_EXIT_PASSED = 0
TEST_EXIT_FAILED = 1

class DATestFramework():

    def __init__(
            self,
            blockchain_node_configs={},
    ):
        self.blockchain_node_configs = blockchain_node_configs

        self.da_encoder_binary = None
        self.da_signer_binary = None
        self.da_server_binary = None
        self.blockchain_binary = None
        self.da_services = []
        self.blockchain_node = None

        tmp_dir = os.path.join(os.path.dirname(__file_path__), "tmp")
        self.__default_blockchain_binary__ = os.path.join(
            tmp_dir, "0gchaind" + binary_ext
        )
        self.__default_da_encoder_binary__ = os.path.join(
            tmp_dir, "da_encoder" + binary_ext
        )
        self.__default_da_signer_binary__ = os.path.join(
            tmp_dir, "da_signer" + binary_ext
        )
        self.__default_da_server_binary__ = os.path.join(
            tmp_dir, "da_server" + binary_ext
        )
        self.tmp_dir = tmp_dir

    def __start_logging(self):
        # Add logger and logging handlers
        self.log = logging.getLogger("TestFramework")
        self.log.setLevel(logging.DEBUG)

        # Create file handler to log all messages
        fh = logging.FileHandler(
            self.options.tmpdir + "/test_framework.log", encoding="utf-8"
        )
        fh.setLevel(logging.DEBUG)

        # Create console handler to log messages to stderr. By default this logs only error messages, but can be configured with --loglevel.
        ch = logging.StreamHandler(sys.stdout)
        # User can provide log level as a number or string (eg DEBUG). loglevel was caught as a string, so try to convert it to an int
        ll = (
            int(self.options.loglevel)
            if self.options.loglevel.isdigit()
            else self.options.loglevel.upper()
        )
        ch.setLevel(ll)

        # Format logs the same as bitcoind's debug.log with microprecision (so log files can be concatenated and sorted)
        formatter = logging.Formatter(
            fmt="%(asctime)s.%(msecs)03d000Z %(name)s (%(levelname)s): %(message)s",
            datefmt="%Y-%m-%dT%H:%M:%S",
        )
        formatter.converter = time.gmtime
        fh.setFormatter(formatter)
        ch.setFormatter(formatter)

        # add the handlers to the logger
        self.log.addHandler(fh)
        self.log.addHandler(ch)

    def add_arguments(self, parser: argparse.ArgumentParser):
        parser.add_argument(
            "--da-encoder-binary",
            dest="da_encoder",
            default=self.__default_da_encoder_binary__,
            type=str,
        )

        parser.add_argument(
            "--da-signer-binary",
            dest="da_signer",
            default=self.__default_da_signer_binary__,
            type=str,
        )
        parser.add_argument(
            "--blockchain-binary",
            dest="blockchain_binary",
            default=self.__default_blockchain_binary__,
            type=str,
        )
        parser.add_argument(
            "--da-server-binary",
            dest="da_server",
            default=self.__default_da_server_binary__,
            type=str,
        )
        parser.add_argument(
            "-l",
            "--loglevel",
            dest="loglevel",
            default="INFO",
            help="log events at this level and higher to the console. Can be set to DEBUG, INFO, WARNING, ERROR or CRITICAL. Passing --loglevel DEBUG will output all logs to console. Note that logs at all levels are always written to the test_framework.log file in the temporary test directory.",
        )
        parser.add_argument(
            "--tmpdir", dest="tmpdir", help="Root directory for datadirs"
        )

        parser.add_argument(
            "--randomseed", dest="random_seed", type=int, help="Set a random seed"
        )

        parser.add_argument("--port-min", dest="port_min", default=11000, type=int)

        parser.add_argument(
            "--pdbonfailure",
            dest="pdbonfailure",
            default=False,
            action="store_true",
            help="Attach a python debugger if test fails",
        )

    def setup_params(self):
        # self.num_blockchain_nodes = 1
        # self.num_nodes = 1
        pass

    def setup_nodes(self):
        self.blockchain_binary = self.options.blockchain_binary
        self.da_encoder_binary = self.options.da_encoder
        self.da_signer_binary = self.options.da_signer
        self.da_server_binary = self.options.da_server

        self.build_binary()

        assert os.path.exists(self.blockchain_binary), f"blockchain binary not found: {self.blockchain_binary}"
        assert os.path.exists(self.da_encoder_binary), f"da encoder binary not found: {self.da_encoder_binary}"
        assert os.path.exists(self.da_signer_binary), f"da signer binary not found: {self.da_signer_binary}"
        assert os.path.exists(self.da_server_binary), f"da server binary not found: {self.da_server_binary}"

        self.setup_blockchain_node()
        self.setup_da_nodes()

    def stop_nodes(self):
        while self.da_services:
            service = self.da_services.pop()
            service.stop()

    def build_binary(self):
        da_root = os.path.join(__file_path__, "..", "..")
        if not os.path.exists(self.blockchain_binary):
            self.build_blockchain(self.blockchain_binary, os.path.join(da_root, "0g-chain"))

        if not os.path.exists(self.da_encoder_binary):
            self.build_da_encoder(self.da_encoder_binary, os.path.join(da_root, "0g-da-encoder"))
        
        if not os.path.exists(self.da_signer_binary):
            self.build_da_signer(self.da_signer_binary, os.path.join(da_root, "0g-da-signer"))

        if not os.path.exists(self.da_server_binary):
            self.build_da_server(self.da_server_binary, os.path.join(da_root, "disperser", "cmd", "combined_server"))

        # da_disperser_cmd = os.path.join(da_root, "disperser", "cmd")
        # if not os.path.exists(self.da_encoder_binary):
        #     self.build_da_node(self.da_encoder_binary, os.path.join(da_disperser_cmd, "encoder"))
        # if not os.path.exists(self.da_batcher_binary):
        #     self.build_da_node(self.da_batcher_binary, os.path.join(da_disperser_cmd, "batcher"))
        # if not os.path.exists(self.da_server_binary):
        #     self.build_da_node(self.da_server_binary, os.path.join(da_disperser_cmd, "apiserver"))
            
        # da_retriever_cmd = os.path.join(da_root, "retriever", "cmd")
        # if not os.path.exists(self.da_retriever_binary):
        #     self.build_da_node(self.da_retriever_binary, da_retriever_cmd)

    def build_blockchain(self, blockchain_binary_path, blockchain_path):            
        origin_path = os.getcwd()
        os.chdir(blockchain_path)
        os.system(f"make build")

        target_path = os.path.join(blockchain_path, "out", "linux",  "0gchaind" + binary_ext)
        script_path = os.path.join(blockchain_path, "localtestnet.sh")
        shutil.copyfile(target_path, blockchain_binary_path)

        sh_path = os.path.join(self.tmp_dir, "localtestnet.sh")
        shutil.copyfile(script_path, sh_path)

        with open(sh_path, 'r+') as file:
            lines = file.readlines()
            file.seek(0)
            file.truncate()
            i = len(lines) - 1
            while len(lines[i].strip()) == 0:
                i -= 1

            file.writelines(lines[:i])

        if not is_windows_platform():
            st = os.stat(blockchain_binary_path)
            os.chmod(blockchain_binary_path, st.st_mode | stat.S_IEXEC)
            os.chmod(sh_path, st.st_mode | stat.S_IEXEC)

        os.chdir(origin_path)

    def build_da_encoder(self, binary_path, src_code_path):
        target_path = os.path.join(src_code_path, "target")
        if os.path.exists(target_path):
            shutil.rmtree(target_path)

        origin_path = os.getcwd()
        os.chdir(src_code_path)
        os.system("cargo build --release --features parallel")

        path = os.path.join(target_path, "release", "server" + binary_ext)
        shutil.copyfile(path, binary_path)

        if not is_windows_platform():
            st = os.stat(binary_path)
            os.chmod(binary_path, st.st_mode | stat.S_IEXEC)

        os.chdir(origin_path)
    
    def build_da_signer(self, binary_path, src_code_path):
        target_path = os.path.join(src_code_path, "target")
        if os.path.exists(target_path):
            shutil.rmtree(target_path)

        origin_path = os.getcwd()
        os.chdir(src_code_path)
        os.system("cargo build --release")

        path = os.path.join(target_path, "release", "server" + binary_ext)
        shutil.copyfile(path, binary_path)

        if not is_windows_platform():
            st = os.stat(binary_path)
            os.chmod(binary_path, st.st_mode | stat.S_IEXEC)
            
        os.chdir(origin_path)

    def build_da_server(self, binary_path, src_code_path):
        os.system(f"go build -o {binary_path} {src_code_path}")

        if not is_windows_platform():
            st = os.stat(binary_path)
            os.chmod(binary_path, st.st_mode | stat.S_IEXEC)

    def build_da_node(self, binary_path, src_code_path):
        os.system(f"go build -o {binary_path} {src_code_path}")

        if not is_windows_platform():
            st = os.stat(binary_path)
            os.chmod(binary_path, st.st_mode | stat.S_IEXEC)

    def setup_blockchain_node(self):
        updated_config = {}

        exe_file = os.path.join(self.root_dir, ".build")
        if not os.path.exists(exe_file):
            os.mkdir(exe_file)
        
        exe_file = os.path.join(exe_file, "0gchaind")
        shutil.copyfile(self.blockchain_binary, exe_file)
        shutil.copystat(self.blockchain_binary, exe_file)

        start_script = os.path.join(self.tmp_dir, "localtestnet.sh")
        log_file = tempfile.NamedTemporaryFile(dir=self.root_dir, delete=False, prefix="init_genesis_", suffix=".log")
        ret = subprocess.run(
            args=["bash", start_script],
            cwd=self.root_dir,
            stdout=log_file,
            stderr=log_file,
        )

        assert ret.returncode == 0, "Failed to init 0gchain genesis, see more details in log file: %s" % log_file.name
        
        rpc_url = "http://0.0.0.0:%s" % blockchain_rpc_port(0)
        self.blockchain_node = BlockchainNode(self.root_dir, rpc_url, self.blockchain_binary, updated_config, self.contract_path, self.log)

        self.blockchain_node.start()

        self.blockchain_node.wait_for_rpc_connection()
        self.blockchain_node.setup_contract()

    def setup_da_nodes(self):
        self.log.info("Start deploy DA services")

        rpc_endpoint =  "http://0.0.0.0:%s" % blockchain_rpc_port(0)
        encoder_endpoint = "0.0.0.0:34001"
        signer_endpoint = "0.0.0.0:34000"

        updated_config = dict(
            params_dir = amt_params,
            grpc_listen_address = encoder_endpoint,
        )
        self.setup_da_node(DAEncoder, self.da_encoder_binary, updated_config)
        updated_config = dict(
            encoder_params_dir = amt_params,
            eth_rpc_endpoint = rpc_endpoint,
            grpc_listen_address = signer_endpoint,
            socket_address = signer_endpoint,
        )
        self.setup_da_node(DASigner, self.da_signer_binary, updated_config)
        
        updated_config = dict(
            blockchain_rpc_endpoint = rpc_endpoint,
            encoder_socket = encoder_endpoint,
            grpc_port = "51001",
        )
        self.setup_da_node(DAServer, self.da_server_binary, updated_config)
        
        self.log.info("All DA service started")

    def setup_da_node(self, clazz, binary, updated_config={}):
        srv = clazz(self.root_dir, binary, updated_config, self.log)
        self.da_services.append(srv)
        srv.setup_config()
        srv.start()
        srv.wait_for_rpc_connection()


    def run_test(self):
        raise NotImplementedError

    def main(self):
        parser = argparse.ArgumentParser(usage="%(prog)s [options]")
        self.add_arguments(parser)
        self.options = parser.parse_args()
        PortMin.n = self.options.port_min

        # Set up temp directory and start logging
        if self.options.tmpdir:
            self.options.tmpdir = os.path.abspath(self.options.tmpdir)
            os.makedirs(self.options.tmpdir, exist_ok=True)
        else:
            self.options.tmpdir = os.getenv(
                "DA_TESTS_LOG_DIR", default=tempfile.mkdtemp(prefix="da_test_")
            )

        self.root_dir = self.options.tmpdir

        self.__start_logging()
        self.log.info("Root dir: %s", self.root_dir)

        self.contract_path = os.path.join(__file_path__, "..", "..", "0g-da-contract")
        assert os.path.exists(self.contract_path), (
            "%s should be exist" % self.contract_path
        )

        if self.options.random_seed is not None:
            random.seed(self.options.random_seed)

        success = TestStatus.FAILED
        try:
            self.setup_params()
            self.setup_nodes()
            self.run_test()
            success = TestStatus.PASSED
        except AssertionError as e:
            self.log.exception("Assertion failed %s", repr(e))
        except KeyboardInterrupt as e:
            self.log.warning("Exiting after keyboard interrupt %s", repr(e))
        except Exception as e:
            self.log.error("Test exception %s %s", repr(e), traceback.format_exc())
            self.log.error(f"Test data are not deleted: {self.root_dir}")

        if success == TestStatus.FAILED and self.options.pdbonfailure:
            print("Testcase failed. Attaching python debugger. Enter ? for help")
            pdb.set_trace()

        if success == TestStatus.PASSED:
            self.log.info("Tests successful")
            exit_code = TEST_EXIT_PASSED
        else:
            self.log.error(
                "Test failed. Test logging available at %s/test_framework.log",
                self.options.tmpdir,
            )
            exit_code = TEST_EXIT_FAILED

        self.stop_nodes()

        if success == TestStatus.PASSED:
            self.log.info("Test success")
            shutil.rmtree(self.root_dir)

        handlers = self.log.handlers[:]
        for handler in handlers:
            self.log.removeHandler(handler)
            handler.close()
        logging.shutdown()

        sys.exit(exit_code)