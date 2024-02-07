package retriever

import (
	"context"
	"fmt"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	pb "github.com/zero-gravity-labs/zerog-data-avail/api/grpc/retriever"
	"github.com/zero-gravity-labs/zerog-data-avail/clients"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/core"
	"github.com/zero-gravity-labs/zerog-storage-client/kv"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

type Server struct {
	pb.UnimplementedRetrieverServer

	config          *Config
	retrievalClient clients.RetrievalClient
	logger          common.Logger
	metrics         *Metrics
	KVNode          *kv.Client
	StreamId        eth_common.Hash
}

func NewServer(
	config *Config,
	logger common.Logger,
	retrievalClient clients.RetrievalClient,
	encoder core.Encoder,
) *Server {
	metrics := NewMetrics(config.MetricsConfig.HTTPPort, logger)

	return &Server{
		config:          config,
		retrievalClient: retrievalClient,
		logger:          logger,
		metrics:         metrics,
		KVNode:          kv.NewClient(node.MustNewClient(config.StorageNodeConfig.KVNodeURL), nil),
		StreamId:        config.StorageNodeConfig.KVStreamId,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.metrics.Start(ctx)
	return nil
}

func (s *Server) fetchBatchInfo(batchHeaderHash []byte) (*core.KVBatchInfo, error) {
	val, err := s.KVNode.GetValue(s.StreamId, batchHeaderHash)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get batch header from kv node")
	}
	batchInfo, err := new(core.KVBatchInfo).Deserialize(val.Data)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to deserialize batch header")
	}
	return batchInfo, nil
}

func (s *Server) RetrieveBlob(ctx context.Context, req *pb.BlobRequest) (*pb.BlobReply, error) {
	s.logger.Info("Received request: ", "BatchHeaderHash", req.GetBatchHeaderHash(), "BlobIndex", req.GetBlobIndex())
	s.metrics.IncrementRetrievalRequestCounter()
	if len(req.GetBatchHeaderHash()) != 32 {
		return nil, fmt.Errorf("got invalid batch header hash")
	}
	var batchHeaderHash [32]byte
	copy(batchHeaderHash[:], req.GetBatchHeaderHash())

	batchInfo, err := s.fetchBatchInfo(req.GetBatchHeaderHash())
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("server fetched batch header: %v", *batchInfo.BatchHeader)

	data, err := s.retrievalClient.RetrieveBlob(
		ctx,
		batchHeaderHash,
		batchInfo.BatchHeader.DataRoot,
		req.GetBlobIndex(),
		uint(batchInfo.BatchHeader.ReferenceBlockNumber),
		batchInfo.BatchHeader.BatchRoot,
		batchInfo.BlobLengths,
	)
	if err != nil {
		return nil, err
	}
	return &pb.BlobReply{
		Data: data,
	}, nil
}
