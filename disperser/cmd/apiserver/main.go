package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zero-gravity-labs/zgda/common"
	"github.com/zero-gravity-labs/zgda/disperser/apiserver"
	"github.com/zero-gravity-labs/zgda/disperser/common/blobstore"

	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zgda/common/aws/dynamodb"
	"github.com/zero-gravity-labs/zgda/common/aws/s3"
	"github.com/zero-gravity-labs/zgda/common/logging"
	"github.com/zero-gravity-labs/zgda/common/ratelimit"
	"github.com/zero-gravity-labs/zgda/common/store"
	"github.com/zero-gravity-labs/zgda/disperser"
	"github.com/zero-gravity-labs/zgda/disperser/cmd/apiserver/flags"
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
	app.Usage = "EigenDA Disperser Server"
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
	var ratelimiter common.RateLimiter

	s3Client, err := s3.NewClient(context.Background(), config.AwsClientConfig, logger)
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
	blobStore = blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)

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
			bucketStore, err = store.NewLocalParamStore[common.RateBucketParams](config.BucketStoreSize)
			if err != nil {
				return err
			}
		}
		ratelimiter = ratelimit.NewRateLimiter(globalParams, bucketStore, config.RatelimiterConfig.Allowlist, logger)
	}

	// TODO: create a separate metrics for batcher
	metrics := disperser.NewMetrics(config.MetricsConfig.HTTPPort, logger)
	server := apiserver.NewDispersalServer(config.ServerConfig, blobStore, logger, metrics, ratelimiter, config.RateConfig)

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Disperser", "socket", httpSocket)
	}

	return server.Start(context.Background())
}
