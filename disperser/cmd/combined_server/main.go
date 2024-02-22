package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/disperser/apiserver"
	"github.com/0glabs/0g-data-avail/disperser/batcher"
	"github.com/0glabs/0g-data-avail/disperser/batcher/dispatcher"
	"github.com/0glabs/0g-data-avail/disperser/batcher/transactor"
	"github.com/0glabs/0g-data-avail/disperser/common/blobstore"
	"github.com/0glabs/0g-data-avail/disperser/common/memorydb"
	"github.com/0glabs/0g-data-avail/disperser/encoder"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/0glabs/0g-data-avail/common/aws/dynamodb"
	"github.com/0glabs/0g-data-avail/common/aws/s3"
	"github.com/0glabs/0g-data-avail/common/geth"
	"github.com/0glabs/0g-data-avail/common/logging"
	"github.com/0glabs/0g-data-avail/common/ratelimit"
	"github.com/0glabs/0g-data-avail/common/store"
	"github.com/0glabs/0g-data-avail/disperser"
	"github.com/0glabs/0g-data-avail/disperser/cmd/combined_server/flags"
	"github.com/urfave/cli"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "combined-server"
	app.Usage = "ZGDA Combined Server"
	app.Description = "Service for disperser server and batcher"

	app.Action = RunCombinedServer
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunDisperserServer(config Config, blobStore disperser.BlobStore, logger common.Logger) error {
	var ratelimiter common.RateLimiter
	if config.EnableRatelimiter {
		globalParams := config.RatelimiterConfig.GlobalRateParams

		var bucketStore common.KVStore[common.RateBucketParams]
		if config.BucketTableName != "" {
			dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
			if err != nil {
				return err
			}
			bucketStore = store.NewDynamoParamStore[common.RateBucketParams](dynamoClient, config.BucketTableName)
		} else {
			var err error
			bucketStore, err = store.NewLocalParamStore[common.RateBucketParams](config.BucketStoreSize)
			if err != nil {
				return err
			}
		}
		ratelimiter = ratelimit.NewRateLimiter(globalParams, bucketStore, config.RatelimiterConfig.Allowlist, logger)
	}

	metrics := disperser.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	var kvClient *kv.Client
	var rpcClient *rpc.Client

	if config.BlobstoreConfig.MetadataHashAsBlobKey {
		kvClient = kv.NewClient(node.MustNewClient(config.StorageNodeConfig.KVNodeURL), nil)
		var err error
		rpcClient, err = rpc.Dial(config.EthClientConfig.RPCURL)
		if err != nil {
			return err
		}
	}
	server := apiserver.NewDispersalServer(config.ServerConfig, blobStore, logger, metrics, ratelimiter, config.RateConfig, config.BlobstoreConfig.MetadataHashAsBlobKey, kvClient, config.StorageNodeConfig.KVStreamId, rpcClient)

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Disperser", "socket", httpSocket)
	}

	return server.Start(context.Background())
}

func RunBatcher(config Config, queue disperser.BlobStore, logger common.Logger) error {
	// transactor
	transactor := transactor.NewTransactor(logger)
	// dispatcher
	dispatcher, err := dispatcher.NewDispatcher(&dispatcher.Config{
		EthClientURL:      config.EthClientConfig.RPCURL,
		PrivateKeyString:  config.EthClientConfig.PrivateKeyString,
		StorageNodeConfig: config.StorageNodeConfig,
	}, transactor, logger)
	if err != nil {
		return err
	}

	// eth clients
	client, err := geth.NewClient(config.EthClientConfig, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}

	rpcClient, err := rpc.Dial(config.EthClientConfig.RPCURL)
	if err != nil {
		return err
	}

	metrics := batcher.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	// encoder
	if len(config.BatcherConfig.EncoderSocket) == 0 {
		return fmt.Errorf("encoder socket must be specified")
	}
	encoderClient, err := encoder.NewEncoderClient(config.BatcherConfig.EncoderSocket, config.TimeoutConfig.EncodingTimeout)
	if err != nil {
		return err
	}

	// confirmer
	confirmer, err := batcher.NewConfirmer(config.EthClientConfig, config.StorageNodeConfig, queue, config.BatcherConfig.MaxNumRetriesPerBlob, config.BatcherConfig.ConfirmerNum, transactor, logger, metrics)
	if err != nil {
		return err
	}

	//finalizer
	finalizer := batcher.NewFinalizer(config.TimeoutConfig.ChainReadTimeout, config.BatcherConfig.FinalizerInterval, queue, client, rpcClient, config.BatcherConfig.MaxNumRetriesPerBlob, logger)

	//batcher
	batcher, err := batcher.NewBatcher(config.BatcherConfig, config.TimeoutConfig, queue, dispatcher, encoderClient, finalizer, confirmer, logger, metrics)
	if err != nil {
		return err
	}

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Batcher", "socket", httpSocket)
	}

	return batcher.Start(context.Background())
}

func RunCombinedServer(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	var blobStore disperser.BlobStore

	if !config.BlobstoreConfig.InMemory {
		s3Client, err := s3.NewClient(config.AwsClientConfig, logger)
		if err != nil {
			return err
		}

		dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
		if err != nil {
			return err
		}

		bucketName := config.BlobstoreConfig.BucketName
		logger.Info("Creating blob store", "bucket", bucketName)
		blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, 0)
		blobStore = blobstore.NewSharedStorage(bucketName, s3Client, config.BlobstoreConfig.MetadataHashAsBlobKey, blobMetadataStore, logger)
	} else {
		config.BlobstoreConfig.MetadataHashAsBlobKey = true
		blobStore = memorydb.NewBlobStore(config.BlobstoreConfig.MemoryDBSize, logger)
	}
	errChan := make(chan error)
	go func() {
		err := RunDisperserServer(config, blobStore, logger)
		errChan <- err
	}()
	go func() {
		err := RunBatcher(config, blobStore, logger)
		errChan <- err
	}()
	err = <-errChan
	return err
}
