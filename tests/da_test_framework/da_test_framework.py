import os
import sys
import shutil
import stat
import argparse

sys.path.append("../0g-storage-kv/tests")

from test_framework.test_framework import TestFramework
from test_framework.blockchain_node import BlockChainNodeType
from utility.kv import MAX_STREAM_ID, to_stream_id
from utility.utils import is_windows_platform
from test_all import build_conflux
from da_test_framework.local_stack import LocalStack
from da_test_framework.da_encoder import DAEncoder
from da_test_framework.da_batcher import DABatcher
from da_test_framework.da_server import DAServer

__file_path__ = os.path.dirname(os.path.realpath(__file__))
binary_ext = ".exe" if is_windows_platform() else ""


class DATestFramework(TestFramework):

    def __init__(
            self,
            blockchain_node_type=BlockChainNodeType.Conflux,
            blockchain_node_configs={},
    ):
        super(DATestFramework, self).__init__(blockchain_node_type, blockchain_node_configs)
        self.localstack_binary = None
        self.da_encoder_binary = None
        self.da_batcher_binary = None
        self.da_server_binary = None
        self.stream_ids = None
        self.da_services = []
        tmp_dir = os.path.join(os.path.dirname(__file_path__), "tmp")
        self.__default_conflux_binary__ = os.path.join(
            tmp_dir, "conflux" + binary_ext
        )
        self.__default_geth_binary__ = os.path.join(
            tmp_dir, "geth" + binary_ext
        )
        self.__default_zgs_node_binary__ = os.path.join(
            tmp_dir, "zgs_node" + binary_ext
        )
        self.__default_zgs_cli_binary__ = os.path.join(
            tmp_dir, "0g-storage-client" + binary_ext
        )
        self.__default_zgs_kv_binary__ = os.path.join(
            tmp_dir, "zgs_kv" + binary_ext
        )
        self.__default_localstack_binary__ = os.path.join(
            tmp_dir, "localstack" + binary_ext
        )
        self.__default_da_encoder_binary__ = os.path.join(
            tmp_dir, "da_encoder" + binary_ext
        )
        self.__default_da_batcher_binary__ = os.path.join(
            tmp_dir, "da_batcher" + binary_ext
        )
        self.__default_da_server_binary__ = os.path.join(
            tmp_dir, "da_server" + binary_ext
        )

    def add_arguments(self, parser: argparse.ArgumentParser):
        super().add_arguments(parser)
        parser.add_argument(
            "--localstack-binary",
            dest="localstack",
            default=self.__default_localstack_binary__,
            type=str,
        )

        parser.add_argument(
            "--da-encoder-binary",
            dest="da_encoder",
            default=self.__default_da_encoder_binary__,
            type=str,
        )

        parser.add_argument(
            "--da-batcher-binary",
            dest="da_batcher",
            default=self.__default_da_batcher_binary__,
            type=str,
        )

        parser.add_argument(
            "--da-server-binary",
            dest="da_server",
            default=self.__default_da_server_binary__,
            type=str,
        )

    def setup_nodes(self):
        self.localstack_binary = self.options.localstack
        self.da_encoder_binary = self.options.da_encoder
        self.da_batcher_binary = self.options.da_batcher
        self.da_server_binary = self.options.da_server

        self.build_binary()

        assert os.path.exists(self.blockchain_binary), f"blockchain binary not found: {self.blockchain_binary}"
        assert os.path.exists(self.zgs_binary), f"zgs binary not found: {self.zgs_binary}"
        assert os.path.exists(self.cli_binary), f"cli binary not found: {self.cli_binary}"
        assert os.path.exists(self.kv_binary), f"kv binary not found: {self.kv_binary}"
        assert os.path.exists(self.localstack_binary), f"localstack binary not found: {self.localstack_binary}"
        assert os.path.exists(self.da_encoder_binary), f"da encoder binary not found: {self.da_encoder_binary}"
        assert os.path.exists(self.da_batcher_binary), f"da batcher binary not found: {self.da_batcher_binary}"
        assert os.path.exists(self.da_server_binary), f"da server binary not found: {self.da_server_binary}"

        super().setup_nodes()
        self.stream_ids = [to_stream_id(i) for i in range(MAX_STREAM_ID)]
        self.stream_ids.reverse()
        super().setup_kv_node(0, self.stream_ids)
        self.setup_da_nodes()

    def stop_nodes(self):
        while self.da_services:
            service = self.da_services.pop()
            service.stop()
        super().stop_nodes()

    def build_binary(self):
        if not os.path.exists(self.blockchain_binary):
            build_conflux(self.blockchain_binary)
        if not os.path.exists(self.zgs_binary) or not os.path.exists(self.cli_binary):
            self.build_zgs_node(self.zgs_binary, self.cli_binary)
        if not os.path.exists(self.kv_binary):
            self.build_zgs_kv(self.kv_binary)

        da_root = os.path.join(__file_path__, "..", "..")
        if not os.path.exists(self.localstack_binary):
            self.build_da_node(self.localstack_binary, os.path.join(da_root, "inabox", "deploy", "cmd"))

        da_disperser_cmd = os.path.join(da_root, "disperser", "cmd")
        if not os.path.exists(self.da_encoder_binary):
            self.build_da_node(self.da_encoder_binary, os.path.join(da_disperser_cmd, "encoder"))
        if not os.path.exists(self.da_batcher_binary):
            self.build_da_node(self.da_batcher_binary, os.path.join(da_disperser_cmd, "batcher"))
        if not os.path.exists(self.da_server_binary):
            self.build_da_node(self.da_server_binary, os.path.join(da_disperser_cmd, "apiserver"))

    def build_zgs_node(self, zgs_node_path, zgs_cli_path):
        zgs_root_path = os.path.join(__file_path__, "..", "..", "0g-storage-kv", "0g-storage-node")
        target_path = os.path.join(zgs_root_path, "target")
        if os.path.exists(target_path):
            shutil.rmtree(target_path)

        origin_path = os.getcwd()
        os.chdir(zgs_root_path)
        os.system("cargo build --release")

        path = os.path.join(target_path, "release", "zgs_node" + binary_ext)
        shutil.copyfile(path, zgs_node_path)
        path = os.path.join(target_path, "0g-storage-client" + binary_ext)
        shutil.copyfile(path, zgs_cli_path)

        if not is_windows_platform():
            st = os.stat(zgs_node_path)
            os.chmod(zgs_node_path, st.st_mode | stat.S_IEXEC)
            st = os.stat(zgs_cli_path)
            os.chmod(zgs_cli_path, st.st_mode | stat.S_IEXEC)

        os.chdir(origin_path)

    def build_zgs_kv(self, kv_path):
        kv_root_path = os.path.join(__file_path__, "..", "..", "0g-storage-kv")
        target_path = os.path.join(kv_root_path, "target")
        if os.path.exists(target_path):
            shutil.rmtree(target_path)

        origin_path = os.getcwd()
        os.chdir(kv_root_path)
        os.system("cargo build --release")

        path = os.path.join(target_path, "release", "zgs_kv" + binary_ext)
        shutil.copyfile(path, kv_path)

        if not is_windows_platform():
            st = os.stat(kv_path)
            os.chmod(kv_path, st.st_mode | stat.S_IEXEC)

        os.chdir(origin_path)

    def build_da_node(self, binary_path, src_code_path):
        os.system(f"go build -o {binary_path} {src_code_path}")

        if not is_windows_platform():
            st = os.stat(binary_path)
            os.chmod(binary_path, st.st_mode | stat.S_IEXEC)

    def setup_da_nodes(self):
        self.log.info("Start deploy DA services")
        self.setup_da_node(LocalStack, self.localstack_binary)
        self.setup_da_node(DAEncoder, self.da_encoder_binary)
        self.setup_da_node(DABatcher, self.da_batcher_binary)
        self.setup_da_node(DAServer, self.da_server_binary)
        self.log.info("All DA service started")

    def setup_da_node(self, clazz, binary, updated_config={}):
        if clazz == DABatcher:
            srv = clazz(self.root_dir, binary, updated_config, self.contract.address(), self.log)
        else:
            srv = clazz(self.root_dir, binary, updated_config, self.log)
        self.da_services.append(srv)
        srv.setup_config()
        srv.start()
        srv.wait_for_rpc_connection()
