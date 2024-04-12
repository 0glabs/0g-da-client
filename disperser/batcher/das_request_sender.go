package batcher

import (
	eth_common "github.com/ethereum/go-ethereum/common"
)

type DASRequest struct {
	StreamId       eth_common.Hash
	BlobHeaderHash [32]byte
	BlobNum        uint
}

type DASRequestSender struct {
}

func NewDASRequestSender() *DASRequestSender {
	return &DASRequestSender{}
}
