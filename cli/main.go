package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/zero-gravity-labs/zerog-data-avail/cli/flags"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws/s3"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "aws-cli"
	app.Usage = "ZGDA CLI"
	app.Description = "Service for S3 Operations"
	app.Flags = flags.Flags
	app.Action = AwsCli
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func AwsCli(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	ctx_bg := context.Background()

	s3Client, err := s3.NewClient(ctx_bg, config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	bucketName := ctx.String(flags.S3BucketNameFlag.Name)
	region := ctx.String(config.AwsClientConfig.Region)
	err = s3Client.CreateBucket(ctx_bg, bucketName, region)

	if err != nil {
		return err
	}

	return nil
}
