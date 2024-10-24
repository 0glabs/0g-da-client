package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0glabs/0g-da-client/common/aws/dynamodb"
	"github.com/0glabs/0g-da-client/common/aws/s3"
	"github.com/0glabs/0g-da-client/common/geth"
	"github.com/0glabs/0g-da-client/common/logging"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/batcher"
	"github.com/0glabs/0g-da-client/disperser/batcher/dispatcher"
	"github.com/0glabs/0g-da-client/disperser/batcher/transactor"
	"github.com/0glabs/0g-da-client/disperser/cmd/batcher/flags"
	"github.com/0glabs/0g-da-client/disperser/common/blobstore"
	"github.com/0glabs/0g-da-client/disperser/contract"
	"github.com/0glabs/0g-da-client/disperser/encoder"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
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
	app.Name = "batcher"
	app.Usage = "ZGDA Batcher"
	app.Description = "Service for creating a batch from queued blobs, distributing coded chunks to nodes, and confirming onchain"

	app.Action = RunBatcher
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunBatcher(ctx *cli.Context) error {
	config := NewConfig(ctx)

	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	daEntranceAddress := eth_common.HexToAddress(config.BatcherConfig.DAEntranceContractAddress)
	daSignersAddress := eth_common.HexToAddress(config.BatcherConfig.DASignersContractAddress)
	daContract, err := contract.NewDAContract(daEntranceAddress, daSignersAddress, config.EthClientConfig.RPCURL, config.EthClientConfig.PrivateKeyString, logger)
	if err != nil {
		return fmt.Errorf("failed to create DAEntrance contract: %w", err)
	}

	// transactor
	transactor := transactor.NewTransactor(config.BatcherConfig.VerifiedCommitRootsTxGasLimit, logger)
	// dispatcher
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

	// blob store
	var queue disperser.BlobStore

	bucketName := config.BlobstoreConfig.BucketName
	s3Client, err := s3.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}
	logger.Info("Initialized S3 client", "bucket", bucketName)

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}
	blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.BlobstoreConfig.TableName, 0)
	queue = blobstore.NewSharedStorage(bucketName, s3Client, config.BlobstoreConfig.MetadataHashAsBlobKey, blobMetadataStore, logger)

	metrics := batcher.NewMetrics(config.MetricsConfig.HTTPPort, logger)

	// Create new store
	kvStore, err := disperser.NewLevelDBStore(config.StorageNodeConfig.KvDbPath+"/chunk", config.StorageNodeConfig.TimeToExpire, logger)
	if err != nil {
		logger.Error("create level db failed")
		return nil
	}

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
	batcher, err := batcher.NewBatcher(config.BatcherConfig,
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

	err = batcher.Start(context.Background())
	if err != nil {
		return err
	}

	return nil

}
