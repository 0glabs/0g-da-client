package main

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws"
	"github.com/zero-gravity-labs/zerog-data-avail/common/geth"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
	"github.com/zero-gravity-labs/zerog-data-avail/common/ratelimit"
	"github.com/zero-gravity-labs/zerog-data-avail/common/storage_node"
	"github.com/zero-gravity-labs/zerog-data-avail/core/encoding"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser/apiserver"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser/batcher"
	server_flags "github.com/zero-gravity-labs/zerog-data-avail/disperser/cmd/apiserver/flags"
	batcher_flags "github.com/zero-gravity-labs/zerog-data-avail/disperser/cmd/batcher/flags"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser/cmd/combined_server/flags"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser/common/blobstore"
)

type Config struct {
	// api server
	AwsClientConfig   aws.ClientConfig
	BlobstoreConfig   blobstore.Config
	ServerConfig      disperser.ServerConfig
	LoggerConfig      logging.Config
	MetricsConfig     disperser.MetricsConfig
	RatelimiterConfig ratelimit.Config
	RateConfig        apiserver.RateConfig
	StorageNodeConfig storage_node.ClientConfig
	EthClientConfig   geth.EthClientConfig
	EnableRatelimiter bool
	BucketTableName   string
	BucketStoreSize   int
	// batcher
	BatcherConfig batcher.Config
	TimeoutConfig batcher.TimeoutConfig
	EncoderConfig encoding.EncoderConfig
}

func NewConfig(ctx *cli.Context) (Config, error) {

	ratelimiterConfig, err := ratelimit.ReadCLIConfig(ctx, server_flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}

	rateConfig, err := apiserver.ReadCLIConfig(ctx)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		// api server
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		ServerConfig: disperser.ServerConfig{
			GrpcPort: ctx.GlobalString(server_flags.GrpcPortFlag.Name),
		},
		EthClientConfig: geth.ReadEthClientConfig(ctx),
		BlobstoreConfig: blobstore.Config{
			BucketName:            ctx.GlobalString(server_flags.S3BucketNameFlag.Name),
			TableName:             ctx.GlobalString(server_flags.DynamoDBTableNameFlag.Name),
			MetadataHashAsBlobKey: ctx.GlobalBool(server_flags.MetadataHashAsBlobKey.Name),
			InMemory:              ctx.GlobalBool(flags.UseMemoryDB.Name),
			MemoryDBSize:          uint64(ctx.GlobalUint(flags.MemoryDBSizeLimit.Name)) * 1024 * 1024,
		},
		LoggerConfig: logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		MetricsConfig: disperser.MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
		RatelimiterConfig: ratelimiterConfig,
		RateConfig:        rateConfig,
		EnableRatelimiter: ctx.GlobalBool(server_flags.EnableRatelimiter.Name),
		BucketTableName:   ctx.GlobalString(server_flags.BucketTableName.Name),
		BucketStoreSize:   ctx.GlobalInt(server_flags.BucketStoreSize.Name),
		StorageNodeConfig: storage_node.ReadClientConfig(ctx, flags.FlagPrefix),
		// batcher
		BatcherConfig: batcher.Config{
			PullInterval:             ctx.GlobalDuration(batcher_flags.PullIntervalFlag.Name),
			FinalizerInterval:        ctx.GlobalDuration(batcher_flags.FinalizerIntervalFlag.Name),
			EncoderSocket:            ctx.GlobalString(batcher_flags.EncoderSocket.Name),
			NumConnections:           ctx.GlobalInt(batcher_flags.NumConnectionsFlag.Name),
			EncodingRequestQueueSize: ctx.GlobalInt(batcher_flags.EncodingRequestQueueSizeFlag.Name),
			BatchSizeMBLimit:         ctx.GlobalUint(batcher_flags.BatchSizeLimitFlag.Name),
			SRSOrder:                 ctx.GlobalInt(batcher_flags.SRSOrderFlag.Name),
			MaxNumRetriesPerBlob:     ctx.GlobalUint(batcher_flags.MaxNumRetriesPerBlobFlag.Name),
			ConfirmerNum:             ctx.GlobalUint(batcher_flags.ConfirmerNumFlag.Name),
		},
		TimeoutConfig: batcher.TimeoutConfig{
			EncodingTimeout:   ctx.GlobalDuration(batcher_flags.EncodingTimeoutFlag.Name),
			ChainReadTimeout:  ctx.GlobalDuration(batcher_flags.ChainReadTimeoutFlag.Name),
			ChainWriteTimeout: ctx.GlobalDuration(batcher_flags.ChainWriteTimeoutFlag.Name),
		},
		EncoderConfig: encoding.ReadCLIConfig(ctx),
	}
	return config, nil
}
