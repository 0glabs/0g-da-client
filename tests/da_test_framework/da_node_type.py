from enum import Enum, unique


@unique
class DANodeType(Enum):
    DA_LOCAL_STACK = 3
    DA_ENCODER = 4
    DA_BATCHER = 5
    DA_SERVER = 6
