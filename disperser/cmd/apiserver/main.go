package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0glabs/0g-da-client/disperser/apiserver"
	"github.com/0glabs/0g-da-client/disperser/common/blobstore"

	"github.com/0glabs/0g-da-client/common/aws/dynamodb"
	"github.com/0glabs/0g-da-client/common/aws/s3"
	"github.com/0glabs/0g-da-client/common/logging"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/cmd/apiserver/flags"
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
	app.Name = "disperser"
	app.Usage = "ZGDA Disperser Server"
	app.Description = "Service for accepting blobs for dispersal"

	app.Action = RunDisperserServer
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func RunDisperserServer(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	var blobStore disperser.BlobStore

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

	// Create new store
	kvStore, err := disperser.NewLevelDBStore(config.StorageNodeConfig.KvDbPath+"/chunk", config.StorageNodeConfig.TimeToExpire, logger)
	if err != nil {
		logger.Error("create level db failed")
		return nil
	}

	// TODO: create a separate metrics for batcher
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
