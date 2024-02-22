package flags

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
	"github.com/zero-gravity-labs/zerog-data-avail/common/ratelimit"
)

const (
	FlagPrefix   = "aws-cli"
	envVarPrefix = "AWS-CLI"
)

var (
	/* Required Flags */
	S3BucketNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "s3-bucket-name"),
		Usage:    "Name of the bucket to store blobs",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "S3_BUCKET_NAME"),
	}
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dynamodb-table-name"),
		Usage:    "Name of the dynamodb table to store blob metadata",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DYNAMODB_TABLE_NAME"),
	}
	BucketTableName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "rate-bucket-table-name"),
		Usage:  "name of the dynamodb table to store rate limiter buckets. If not provided, a local store will be used",
		Value:  "",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "RATE_BUCKET_TABLE_NAME"),
	}
)

var requiredFlags = []cli.Flag{
	S3BucketNameFlag,
	DynamoDBTableNameFlag,
	BucketTableName,
}

var optionalFlags = []cli.Flag{
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, logging.CLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, ratelimit.RatelimiterCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
}
