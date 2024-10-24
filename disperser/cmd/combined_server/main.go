package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/disperser/apiserver"
	"github.com/0glabs/0g-da-client/disperser/batcher"
	"github.com/0glabs/0g-da-client/disperser/batcher/dispatcher"
	"github.com/0glabs/0g-da-client/disperser/batcher/transactor"
	"github.com/0glabs/0g-da-client/disperser/common/blobstore"
	"github.com/0glabs/0g-da-client/disperser/common/memorydb"
	"github.com/0glabs/0g-da-client/disperser/contract"
	"github.com/0glabs/0g-da-client/disperser/encoder"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/0glabs/0g-da-client/common/aws/dynamodb"
	"github.com/0glabs/0g-da-client/common/aws/s3"
	"github.com/0glabs/0g-da-client/common/geth"
	"github.com/0glabs/0g-da-client/common/logging"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/cmd/combined_server/flags"
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

func RunDisperserServer(config Config, blobStore disperser.BlobStore, logger common.Logger, kvStore *disperser.Store) error {
	metrics := disperser.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	server := apiserver.NewDispersalServer(config.ServerConfig, blobStore, logger, metrics, config.RatelimiterConfig, config.EnableRatelimiter, config.RateConfig, config.BlobstoreConfig.MetadataHashAsBlobKey, kvStore, config.RetrieverAddr)

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Disperser", "socket", httpSocket)
	}

	return server.Start(context.Background())
}

func RunBatcher(config Config, queue disperser.BlobStore, logger common.Logger, kvStore *disperser.Store) error {
	// transactor
	transactor := transactor.NewTransactor(config.BatcherConfig.VerifiedCommitRootsTxGasLimit, logger)
	// dispatcher
	daEntranceAddress := eth_common.HexToAddress(config.BatcherConfig.DAEntranceContractAddress)
	daSignersAddress := eth_common.HexToAddress(config.BatcherConfig.DASignersContractAddress)
	daContract, err := contract.NewDAContract(daEntranceAddress, daSignersAddress, config.EthClientConfig.RPCURL, config.EthClientConfig.PrivateKeyString, logger)
	if err != nil {
		return fmt.Errorf("failed to create DAEntrance contract: %w", err)
	}

	dispatcher, err := dispatcher.NewDispatcher(transactor, daContract, logger)
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
	confirmer, err := batcher.NewConfirmer(config.EthClientConfig, config.BatcherConfig, queue, daContract, logger, metrics)
	if err != nil {
		return err
	}

	blobKeyCache := disperser.BlobKeyCache{
		Key:   make(map[[32]byte]bool),
		Epoch: 0,
	}
	iter := kvStore.MetadataIterator(context.Background())
	for iter.Next() {
		metadata, err := new(disperser.BlobRetrieveMetadata).Deserialize(iter.Value())
		if err != nil {
			return err
		}

		blobKeyCache.Add(metadata.Hash(), metadata.Epoch)
	}
	iter.Release()

	//finalizer
	finalizer := batcher.NewFinalizer(config.TimeoutConfig.ChainReadTimeout, config.BatcherConfig, queue, client, rpcClient, logger, kvStore, &blobKeyCache)

	//batcher
	batcher, err := batcher.NewBatcher(
		config.BatcherConfig,
		config.TimeoutConfig,
		config.EthClientConfig,
		queue,
		dispatcher,
		encoderClient,
		finalizer,
		confirmer,
		daContract,
		logger,
		metrics,
		&blobKeyCache)
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

	// Create new store
	kvStore, err := disperser.NewLevelDBStore(config.StorageNodeConfig.KvDbPath+"/chunk", config.StorageNodeConfig.TimeToExpire, logger)
	if err != nil {
		logger.Error("create level db failed")
		return nil
	}

	errChan := make(chan error)
	go func() {
		err := RunDisperserServer(config, blobStore, logger, kvStore)
		errChan <- err
	}()
	go func() {
		err := RunBatcher(config, blobStore, logger, kvStore)
		errChan <- err
	}()
	err = <-errChan
	return err
}
