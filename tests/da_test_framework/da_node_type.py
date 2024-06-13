from enum import Enum, unique


@unique
class DANodeType(Enum):
    DA_LOCAL_STACK = 3
    DA_ENCODER = 4
    DA_BATCHER = 5
    DA_SERVER = 6
    DA_RETRIEVER = 7
    DA_SIGNER = 8
