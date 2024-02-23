package flags

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
)

const (
	FlagPrefix   = "aws-cli"
	envVarPrefix = "AWS_CLI"
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
		Name:     common.PrefixFlag(FlagPrefix, "table-name"),
		Usage:    "Name of the dynamodb table",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "TABLE_NAME"),
	}
)

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(logging.CLIFlags(envVarPrefix, FlagPrefix), aws.ClientFlags(envVarPrefix, FlagPrefix)...)
}
