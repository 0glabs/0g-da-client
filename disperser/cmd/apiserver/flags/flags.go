package flags

import (
	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/common/aws"
	"github.com/0glabs/0g-da-client/common/logging"
	"github.com/0glabs/0g-da-client/common/ratelimit"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "disperser-server"
	EnvVarPrefix = "DISPERSER_SERVER"
)

var (
	/* Required Flags */
	S3BucketNameFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "s3-bucket-name"),
		Usage:  "Name of the bucket to store blobs",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "S3_BUCKET_NAME"),
	}
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "dynamodb-table-name"),
		Usage:  "Name of the dynamodb table to store blob metadata",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "DYNAMODB_TABLE_NAME"),
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-port"),
		Usage:    "Port at which disperser listens for grpc calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "GRPC_PORT"),
	}
	/* Optional Flags*/
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "METRICS_HTTP_PORT"),
	}
	EnableMetrics = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_METRICS"),
	}
	EnableRatelimiter = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "enable-ratelimiter"),
		Usage:  "enable rate limiter",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "ENABLE_RATELIMITER"),
	}
	BucketTableName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "rate-bucket-table-name"),
		Usage:  "name of the dynamodb table to store rate limiter buckets. If not provided, a local store will be used",
		Value:  "",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "RATE_BUCKET_TABLE_NAME"),
	}
	BucketStoreSize = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "rate-bucket-store-size"),
		Usage:    "size (max number of entries) of the local store to use for rate limiting buckets",
		Value:    100_000,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "RATE_BUCKET_STORE_SIZE"),
		Required: false,
	}
	MetadataHashAsBlobKey = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "metadata-hash-as-blob-key"),
		Usage:  "use metadata hash as blob key",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "METADATA_HASH_AS_BLOB_KEY"),
	}
	RetrieverAddrName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "retriever-address"),
		Usage:  "address of retriever",
		Value:  "0.0.0.0:34005",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "RETRIEVER-ADDRESS"),
	}
)

var RequiredFlags = []cli.Flag{
	S3BucketNameFlag,
	DynamoDBTableNameFlag,
	GrpcPortFlag,
	BucketTableName,
}

var OptionalFlags = []cli.Flag{
	MetricsHTTPPort,
	EnableMetrics,
	EnableRatelimiter,
	BucketStoreSize,
	MetadataHashAsBlobKey,
	RetrieverAddrName,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(RequiredFlags, OptionalFlags...)
	Flags = append(Flags, logging.CLIFlags(EnvVarPrefix, FlagPrefix)...)
	Flags = append(Flags, ratelimit.RatelimiterCLIFlags(EnvVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(EnvVarPrefix, FlagPrefix)...)
}
