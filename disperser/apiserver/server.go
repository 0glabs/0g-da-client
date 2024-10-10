package apiserver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	pb "github.com/0glabs/0g-da-client/api/grpc/disperser"
	"github.com/0glabs/0g-da-client/common"
	healthcheck "github.com/0glabs/0g-da-client/common/healthcheck"
	"github.com/0glabs/0g-da-client/common/ratelimit"
	"github.com/0glabs/0g-da-client/core"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/api/grpc/retriever"
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

	rateConfig RateConfig

	metrics *disperser.Metrics

	metadataHashAsBlobKey bool
	kvStore               *disperser.Store

	logger common.Logger

	retrieverAddr string

	writeRateLimiterManager *ClientRateLimiterManager
	readRateLimiterManager  *ClientRateLimiterManager
}

// NewServer creates a new Server struct with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewDispersalServer(
	config disperser.ServerConfig,
	store disperser.BlobStore,
	logger common.Logger,
	metrics *disperser.Metrics,
	ratelimiterConfig ratelimit.Config,
	enableRatelimiter bool,
	rateConfig RateConfig,
	metadataHashAsBlobKey bool,
	kvStore *disperser.Store,
	retrieverAddr string,
) *DispersalServer {

	return &DispersalServer{
		config:                config,
		blobStore:             store,
		metrics:               metrics,
		logger:                logger,
		rateConfig:            rateConfig,
		mu:                    &sync.RWMutex{},
		metadataHashAsBlobKey: metadataHashAsBlobKey,
		kvStore:               kvStore,
		retrieverAddr:         retrieverAddr,

		writeRateLimiterManager: NewClientRateLimiterManager(enableRatelimiter, ratelimiterConfig.MaxWriteRequestPerMinute, ratelimiterConfig.Allowlist),
		readRateLimiterManager:  NewClientRateLimiterManager(enableRatelimiter, ratelimiterConfig.MaxReadRequestPerMinute, ratelimiterConfig.Allowlist),
	}
}

func (s *DispersalServer) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		s.metrics.ObserveLatency("DisperseBlob", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	blobSize := len(req.GetData())
	// The blob size in bytes must be in range [1, maxBlobSize].
	if blobSize > core.MaxBlobSize-4 {
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

	s.logger.Debug("[apiserver] received a new blob request", "origin", origin)

	limiter := s.writeRateLimiterManager.GetRateLimiter(origin)
	if limiter != nil && !limiter.Allow() {
		s.logger.Debug("[apiserver] rate limit exceeded for disperse blob", "client", origin)
		return nil, fmt.Errorf("request ratelimited")
	}

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
			s.logger.Warn("get metadata from kv", "error", err)
			return nil, fmt.Errorf("failed to retrieve blob status: disperse failed or request may not exist")
		}
		if metadataFromKV != nil {
			// metadata = metadataInKV
			metadata = &disperser.BlobMetadata{
				BlobStatus: disperser.Finalized,
				ConfirmationInfo: &disperser.ConfirmationInfo{
					DataRoot: metadataFromKV.DataRoot,
					Epoch:    metadataFromKV.Epoch,
					QuorumId: metadataFromKV.QuorumId,
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

		return &pb.BlobStatusReply{
			Status: getResponseStatus(metadata.BlobStatus),
			Info: &pb.BlobInfo{
				BlobHeader: &pb.BlobHeader{
					StorageRoot: confirmationInfo.DataRoot,
					Epoch:       confirmationInfo.Epoch,
					QuorumId:    confirmationInfo.QuorumId,
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
	if len(req.StorageRoot) != 32 {
		return nil, fmt.Errorf("parameter StorageRoot len is not accepted")
	}

	origin, err := common.GetClientAddress(ctx, s.rateConfig.ClientIPHeader, 2, true)
	if err != nil {
		s.metrics.HandleFailedRequest(0, "RetrieveBlob")
		return nil, err
	}

	limiter := s.readRateLimiterManager.GetRateLimiter(origin)
	if limiter != nil && !limiter.Allow() {
		s.logger.Debug("[apiserver] rate limit exceeded for retrieve blob", "client", origin)
		return nil, fmt.Errorf("request rate limited")
	}

	metaData := disperser.BlobRetrieveMetadata{
		DataRoot: req.StorageRoot,
		Epoch:    req.Epoch,
		QuorumId: req.QuorumId,
	}
	blobKey, err := metaData.Serialize()
	if err != nil {
		s.logger.Error("[apiserver] failed to serialize metadata")
	} else {
		data, err := s.kvStore.GetBlob(ctx, blobKey)
		if err != nil {
			s.logger.Error("[apiserver] failed to get blob for key", "blobKey", blobKey)
		} else {
			s.metrics.HandleSuccessfulRequest(len(data), "RetrieveBlob")
			return &pb.RetrieveBlobReply{
				Data: data,
			}, nil
		}
	}

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
	data := req.GetData()

	blob := &core.Blob{
		Data: data,
	}

	return blob
}

type TinyRateLimiter struct {
	ClientID    string // Identifier for the client (e.g., IP address, user ID)
	MaxRequests int
	LastRequest time.Time  // Time of the last request
	mu          sync.Mutex // Mutex for thread safety
}

func (rl *TinyRateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.LastRequest)
	if elapsed < time.Minute/time.Duration(rl.MaxRequests) {
		// Too soon since the last request
		return false
	}

	// Allow the request and update the last request time
	rl.LastRequest = now
	return true
}

type ClientRateLimiterManager struct {
	clients           map[string]*TinyRateLimiter // Map of client ID to RateLimiter
	maxRequests       int
	allowlist         []string
	EnableRatelimiter bool
	mu                sync.Mutex // Mutex for thread safety
}

func NewClientRateLimiterManager(enableRatelimiter bool, maxRequests int, allowlist []string) *ClientRateLimiterManager {
	return &ClientRateLimiterManager{
		clients:           make(map[string]*TinyRateLimiter),
		maxRequests:       maxRequests,
		allowlist:         allowlist,
		EnableRatelimiter: enableRatelimiter,
	}
}

func (m *ClientRateLimiterManager) GetRateLimiter(clientID string) *TinyRateLimiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.EnableRatelimiter {
		return nil
	}

	for _, id := range m.allowlist {
		if strings.Contains(clientID, id) {
			return nil
		}
	}

	limiter, ok := m.clients[clientID]
	if !ok {
		limiter = &TinyRateLimiter{
			ClientID:    clientID,
			MaxRequests: m.maxRequests,
		}
		m.clients[clientID] = limiter
	}

	return limiter
}
