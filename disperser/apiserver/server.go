package apiserver

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/0glabs/0g-da-client/api/grpc/disperser"
	"github.com/0glabs/0g-da-client/common"
	healthcheck "github.com/0glabs/0g-da-client/common/healthcheck"
	"github.com/0glabs/0g-da-client/core"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/api/grpc/retriever"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var errSystemRateLimit = fmt.Errorf("request ratelimited: system limit")
var errAccountRateLimit = fmt.Errorf("request ratelimited: account limit")

const systemAccountKey = "system"

type DispersalServer struct {
	pb.UnimplementedDisperserServer
	mu *sync.RWMutex

	config disperser.ServerConfig

	blobStore disperser.BlobStore

	rateConfig  RateConfig
	ratelimiter common.RateLimiter

	metrics *disperser.Metrics

	metadataHashAsBlobKey bool
	kvStore               *disperser.Store

	rpcClient *rpc.Client

	logger common.Logger

	retrieverAddr string
}

// NewServer creates a new Server struct with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewDispersalServer(
	config disperser.ServerConfig,
	store disperser.BlobStore,
	logger common.Logger,
	metrics *disperser.Metrics,
	ratelimiter common.RateLimiter,
	rateConfig RateConfig,
	metadataHashAsBlobKey bool,
	rpcClient *rpc.Client,
	kvStore *disperser.Store,
	retrieverAddr string,
) *DispersalServer {

	return &DispersalServer{
		config:                config,
		blobStore:             store,
		metrics:               metrics,
		logger:                logger,
		ratelimiter:           ratelimiter,
		rateConfig:            rateConfig,
		mu:                    &sync.RWMutex{},
		metadataHashAsBlobKey: metadataHashAsBlobKey,
		kvStore:               kvStore,
		rpcClient:             rpcClient,
		retrieverAddr:         retrieverAddr,
	}
}

func (s *DispersalServer) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("DisperseBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	securityParams := req.GetSecurityParams()

	blobSize := len(req.GetData())
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > core.MaxBlobSize {
		return nil, fmt.Errorf("blob size cannot exceed %v KiB", core.MaxBlobSize/1024)
	}
	if blobSize == 0 {
		return nil, fmt.Errorf("blob size must be greater than 0")
	}

	blob := getBlobFromRequest(req)

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		s.metrics.HandleFailedRequest(blobSize, "DisperseBlob")
		return nil, err
	}

	s.logger.Debug("[apiserver] received a new blob request", "origin", origin, "securityParams", securityParams)

	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err := s.blobStore.StoreBlob(ctx, blob, requestedAt)
	if err != nil {
		s.metrics.HandleFailedRequest(blobSize, "DisperseBlob")
		return nil, err
	}

	s.metrics.HandleSuccessfulRequest(blobSize, "DisperseBlob")

	s.logger.Info("[apiserver] received a new blob: ", "key", metadataKey.String())
	return &pb.DisperseBlobReply{
		Result:    pb.BlobStatus_PROCESSING,
		RequestId: []byte(metadataKey.String()),
	}, nil
}

func (s *DispersalServer) getMetadataFromKv(ctx context.Context, key []byte) (*disperser.BlobRetrieveMetadata, error) {
	val, err := s.kvStore.GetMetadata(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob metadata from kv node: %v", err)
	}
	if len(val) == 0 {
		return nil, nil
	}
	metadata, err := new(disperser.BlobRetrieveMetadata).Deserialize(val)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize blob metadata: %v", err)
	}
	return metadata, nil
}

func (s *DispersalServer) GetBlobStatus(ctx context.Context, req *pb.BlobStatusRequest) (*pb.BlobStatusReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("GetBlobStatus", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	requestID := req.GetRequestId()
	if len(requestID) == 0 {
		return nil, fmt.Errorf("invalid request: request_id must not be empty")
	}

	s.logger.Info("[apiserver] received a new blob status request", "requestID", string(requestID))
	metadataKey, err := disperser.ParseBlobKey(string(requestID))
	if err != nil {
		return nil, err
	}

	metadata, err := s.blobStore.GetBlobMetadata(ctx, metadataKey)
	if err != nil && !s.metadataHashAsBlobKey {
		return nil, err
	}
	if (metadata == nil || metadata.GetBlobKey().String() != string(requestID)) && s.metadataHashAsBlobKey {
		// check on kv
		metadataFromKV, err := s.getMetadataFromKv(ctx, requestID)
		if err != nil {
			s.logger.Warn("get metadata from kv", err)
		}
		if metadataFromKV != nil {
			// metadata = metadataInKV
			metadata = &disperser.BlobMetadata{
				BlobStatus: disperser.Processing,
				ConfirmationInfo: &disperser.ConfirmationInfo{
					DataRoot:                metadataFromKV.DataRoot,
					Epoch:                   metadataFromKV.Epoch,
					QuorumId:                metadataFromKV.QuorumId,
					ConfirmationBlockNumber: metadataFromKV.BlockNumber,
				},
			}
		} else {
			// behavior align with aws dynamodb
			metadata = &disperser.BlobMetadata{
				BlobStatus: disperser.Processing,
			}
		}
	}

	isConfirmed, err := metadata.IsConfirmed()
	if err != nil {
		return nil, err
	}

	s.logger.Debug("[apiserver] isConfirmed", "metadata", metadata, "isConfirmed", isConfirmed)
	if isConfirmed {
		confirmationInfo := metadata.ConfirmationInfo

		commitmentRoot := confirmationInfo.CommitmentRoot
		dataLength := uint32(confirmationInfo.Length)
		quorumInfos := confirmationInfo.BlobQuorumInfos
		blobQuorumParams := make([]*pb.BlobQuorumParam, len(quorumInfos))
		quorumNumbers := make([]byte, len(quorumInfos))
		quorumPercentSigned := make([]byte, len(quorumInfos))
		quorumIndexes := make([]byte, len(quorumInfos))
		for i, quorumInfo := range quorumInfos {
			blobQuorumParams[i] = &pb.BlobQuorumParam{
				QuorumNumber:                 uint32(quorumInfo.QuorumID),
				AdversaryThresholdPercentage: uint32(quorumInfo.AdversaryThreshold),
				QuorumThresholdPercentage:    uint32(quorumInfo.QuorumThreshold),
				ChunkLength:                  uint32(quorumInfo.ChunkLength),
			}
			quorumNumbers[i] = quorumInfo.QuorumID
			quorumPercentSigned[i] = confirmationInfo.QuorumResults[quorumInfo.QuorumID].PercentSigned
			quorumIndexes[i] = byte(i)
		}

		return &pb.BlobStatusReply{
			Status: getResponseStatus(metadata.BlobStatus),
			Info: &pb.BlobInfo{
				BlobHeader: &pb.BlobHeader{
					CommitmentRoot:   commitmentRoot,
					DataLength:       dataLength,
					BlobQuorumParams: blobQuorumParams,
					DataRoot:         confirmationInfo.DataRoot,
					Epoch:            confirmationInfo.Epoch,
					QuorumId:         confirmationInfo.QuorumId,
				},
				BlobVerificationProof: &pb.BlobVerificationProof{
					BatchId:   confirmationInfo.BatchID,
					BlobIndex: confirmationInfo.BlobIndex,
					BatchMetadata: &pb.BatchMetadata{
						BatchHeader: &pb.BatchHeader{
							BatchRoot:               confirmationInfo.BatchRoot,
							QuorumNumbers:           quorumNumbers,
							QuorumSignedPercentages: quorumPercentSigned,
							ReferenceBlockNumber:    confirmationInfo.ReferenceBlockNumber,
						},
						SignatoryRecordHash:     confirmationInfo.SignatoryRecordHash[:],
						Fee:                     confirmationInfo.Fee,
						ConfirmationBlockNumber: confirmationInfo.ConfirmationBlockNumber,
						BatchHeaderHash:         confirmationInfo.BatchHeaderHash[:],
					},
					InclusionProof: confirmationInfo.BlobInclusionProof,
					// ref: api/proto/disperser/disperser.proto:BlobVerificationProof.quorum_indexes
					QuorumIndexes: quorumIndexes,
				},
			},
		}, nil
	}

	return &pb.BlobStatusReply{
		Status: getResponseStatus(metadata.BlobStatus),
		Info:   &pb.BlobInfo{},
	}, nil
}

func (s *DispersalServer) RetrieveBlob(ctx context.Context, req *pb.RetrieveBlobRequest) (*pb.RetrieveBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("RetrieveBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	s.logger.Info("[apiserver] received a new blob retrieval request", "blob storage root", req.StorageRoot, "blob epoch", req.Epoch, "quorum id", req.QuorumId)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctxWithTimeout,
		s.retrieverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial retriever: %w", err)
	}
	defer conn.Close()

	client := retriever.NewRetrieverClient(conn)
	reply, err := client.RetrieveBlob(ctx, &retriever.BlobRequest{
		StorageRoot: req.StorageRoot,
		Epoch:       req.Epoch,
		QuorumId:    req.QuorumId,
	})

	data := reply.GetData()
	if err != nil {
		s.logger.Error("Failed to retrieve blob", "err", err)
		s.metrics.HandleFailedRequest(len(data), "RetrieveBlob")

		return nil, err
	}

	s.metrics.HandleSuccessfulRequest(len(data), "RetrieveBlob")

	return &pb.RetrieveBlobReply{
		Data: data,
	}, nil
}

func (s *DispersalServer) Start(ctx context.Context) error {
	s.logger.Trace("Entering Start function...")
	defer s.logger.Trace("Exiting Start function...")

	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterDisperserServer(gs, s)

	// Register Server for Health Checks
	healthcheck.RegisterHealthServer(gs)

	s.logger.Info("[apiserver] port", s.config.GrpcPort, "address", listener.Addr().String(), "GRPC Listening")
	if err := gs.Serve(listener); err != nil {
		return fmt.Errorf("could not start GRPC server")
	}

	return nil
}

func getResponseStatus(status disperser.BlobStatus) pb.BlobStatus {
	switch status {
	case disperser.Processing:
		return pb.BlobStatus_PROCESSING
	case disperser.Confirmed:
		return pb.BlobStatus_CONFIRMED
	case disperser.Failed:
		return pb.BlobStatus_FAILED
	case disperser.Finalized:
		return pb.BlobStatus_FINALIZED
	case disperser.InsufficientSignatures:
		return pb.BlobStatus_INSUFFICIENT_SIGNATURES
	default:
		return pb.BlobStatus_UNKNOWN
	}
}

func getBlobFromRequest(req *pb.DisperseBlobRequest) *core.Blob {
	params := make([]*core.SecurityParam, len(req.SecurityParams))

	for i, param := range req.GetSecurityParams() {
		params[i] = &core.SecurityParam{
			QuorumID:           core.QuorumID(param.QuorumId),
			AdversaryThreshold: uint8(param.AdversaryThreshold),
			QuorumThreshold:    uint8(param.QuorumThreshold),
		}
	}

	data := req.GetData()

	blob := &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: params,
		},
		Data: data,
	}

	return blob
}
