package flags

import (
	"time"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/common/aws"
	"github.com/0glabs/0g-data-avail/common/geth"
	"github.com/0glabs/0g-data-avail/common/logging"
	"github.com/0glabs/0g-data-avail/common/storage_node"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "batcher"
	EnvVarPrefix = "BATCHER"
)

var (
	/* Required Flags */
	S3BucketNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "s3-bucket-name"),
		Usage:    "Name of the bucket to store blobs",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "S3_BUCKET_NAME"),
	}
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dynamodb-table-name"),
		Usage:    "Name of the dynamodb table to store blob metadata",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DYNAMODB_TABLE_NAME"),
	}
	PullIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "pull-interval"),
		Usage:    "Interval at which to pull from the queue",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "PULL_INTERVAL"),
	}
	EncoderSocket = cli.StringFlag{
		Name:     "encoder-socket",
		Usage:    "the http ip:port which the distributed encoder server is listening",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENCODER_ADDRESS"),
	}
	EnableMetrics = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_METRICS"),
	}
	BatchSizeLimitFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "batch-size-limit"),
		Usage:    "the maximum batch size in MiB",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "BATCH_SIZE_LIMIT"),
	}
	/* Optional Flags*/
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "METRICS_HTTP_PORT"),
	}
	EncodingTimeoutFlag = cli.DurationFlag{
		Name:     "encoding-timeout",
		Usage:    "connection timeout from grpc call to encoder",
		Required: false,
		Value:    30 * time.Second,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENCODING_TIMEOUT"),
	}
	ChainReadTimeoutFlag = cli.DurationFlag{
		Name:     "chain-read-timeout",
		Usage:    "connection timeout to read from chain",
		Required: false,
		Value:    5 * time.Second,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "CHAIN_READ_TIMEOUT"),
	}
	ChainWriteTimeoutFlag = cli.DurationFlag{
		Name:     "chain-write-timeout",
		Usage:    "connection timeout to write to chain",
		Required: false,
		Value:    90 * time.Second,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "CHAIN_WRITE_TIMEOUT"),
	}
	NumConnectionsFlag = cli.IntFlag{
		Name:     "num-connections",
		Usage:    "maximum number of connections to encoders (defaults to 256)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "NUM_CONNECTIONS"),
		Value:    256,
	}
	SigningTimeoutFlag = cli.DurationFlag{
		Name:     "signing-timeout",
		Usage:    "connection timeout from grpc call to signer",
		Required: false,
		Value:    30 * time.Second,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "SIGNING_TIMEOUT"),
	}
	FinalizerIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "finalizer-interval"),
		Usage:    "Interval at which to check for finalized blobs",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "FINALIZER_INTERVAL"),
		Value:    6 * time.Minute,
	}
	EncodingRequestQueueSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoding-request-queue-size"),
		Usage:    "Size of the encoding request queue",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENCODING_REQUEST_QUEUE_SIZE"),
		Value:    500,
	}
	MaxNumRetriesPerBlobFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-num-retries-per-blob"),
		Usage:    "Maximum number of retries to process a blob before marking the blob as FAILED",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "MAX_NUM_RETRIES_PER_BLOB"),
		Value:    2,
	}
	ConfirmerNumFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "confirmer-num"),
		Usage:    "Number of confirmer go routines",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "CONFIRMER_NUM"),
		Value:    1,
	}
	DAEntranceContractAddressFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "da-entrance-contract"),
		Usage:    "DAEntrance contract address",
		Required: false,
		Value:    "0x0000000000000000000000000000000000000000",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DAENTRANCE_CONTRACT_ADDRESS"),
	}
	DASignersContractAddressFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "da-signers-contract"),
		Usage:    "DASigners contract address",
		Required: false,
		Value:    "0x0000000000000000000000000000000000000000",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "DASIGNERS_CONTRACT_ADDRESS"),
	}
	EncodingIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoding-interval"),
		Usage:    "Interval for encoding loop",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENCODING_INTERVAL"),
		Value:    10 * time.Second,
	}
	SigningIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signing-interval"),
		Usage:    "Interval for signing loop",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "SIGNING_INTERVAL"),
		Value:    10 * time.Second,
	}
	MaxNumRetriesForSignFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-num-retries-for-sign"),
		Usage:    "Maximum number of retries to sign a blob before marking the blob as FAILED",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "MAX_NUM_RETRIES_FOR_SIGN"),
		Value:    1,
	}
	FinalizedBlockCountFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "finalized-block-count"),
		Usage:    "Number of latest block before finalized",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "FINALIZED_BLOCK_COUNT"),
		Value:    1,
	}
	ExpirationPollIntervalSecFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "expiration-poll-interval"),
		Usage:    "How often (in second) to poll status and expire outdated blobs",
		Required: false,
		Value:    "180",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "EXPIRATION_POLL_INTERVAL"),
	}
	SignedPullIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signed-pull-interval"),
		Usage:    "Interval at which to pull from the signed queue",
		Required: true,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "SIGNED_PULL_INTERVAL"),
	}
	VerifiedCommitRootsTxGasLimitFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "verified-commit-roots-tx-gas-limit"),
		Usage:    "tx gas limit for VerifiedCommitRootsTx",
		Required: false,
		Value:    10000000,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "VERIFIED_COMMIT_ROOTS_TX_GAS_LIMIT"),
	}

	// This flag is available so that we can manually adjust the number of chunks if desired for testing purposes or for other reasons.
	// For instance, we may want to increase the number of chunks / reduce the chunk size to reduce the amount of data that needs to be
	// downloaded by light clients for DAS.
	TargetNumChunksFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "target-num-chunks"),
		Usage:    "Target number of chunks per blob. If set to zero, the number of chunks will be calculated based on the ratio of the total stake to the minimum stake",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "TARGET_NUM_CHUNKS"),
		Value:    0,
	}
	MetadataHashAsBlobKey = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "metadata-hash-as-blob-key"),
		Usage:  "use metadata hash as blob key",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "METADATA_HASH_AS_BLOB_KEY"),
	}
)

var RequiredFlags = []cli.Flag{
	S3BucketNameFlag,
	DynamoDBTableNameFlag,
	PullIntervalFlag,
	EncoderSocket,
	EnableMetrics,
	BatchSizeLimitFlag,
	SignedPullIntervalFlag,
}

var OptionalFlags = []cli.Flag{
	MetricsHTTPPort,
	EncodingTimeoutFlag,
	ChainReadTimeoutFlag,
	ChainWriteTimeoutFlag,
	NumConnectionsFlag,
	FinalizerIntervalFlag,
	EncodingRequestQueueSizeFlag,
	MaxNumRetriesPerBlobFlag,
	ConfirmerNumFlag,
	SigningTimeoutFlag,
	DAEntranceContractAddressFlag,
	DASignersContractAddressFlag,
	EncodingIntervalFlag,
	SigningIntervalFlag,
	MaxNumRetriesForSignFlag,
	FinalizedBlockCountFlag,
	ExpirationPollIntervalSecFlag,
	TargetNumChunksFlag,
	MetadataHashAsBlobKey,
	VerifiedCommitRootsTxGasLimitFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(RequiredFlags, OptionalFlags...)
	Flags = append(Flags, geth.EthClientFlags(EnvVarPrefix)...)
	Flags = append(Flags, logging.CLIFlags(EnvVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(EnvVarPrefix, FlagPrefix)...)
	Flags = append(Flags, storage_node.ClientFlags(EnvVarPrefix, FlagPrefix)...)
}
