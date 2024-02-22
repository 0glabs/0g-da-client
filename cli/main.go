package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/zero-gravity-labs/zerog-data-avail/cli/flags"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws/dynamodb"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws/s3"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
	"github.com/zero-gravity-labs/zerog-data-avail/common/store"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser/common/blobstore"
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
	app.Commands = []cli.Command{
		{
		  Name:    "bucket",
		  Usage:   "bucket operation in s3",
		  Subcommands: []cli.Command{
			{
			  Name:  "create",
			  Aliases: []string{"c"},
			  Usage: "create a new bucket",
			  Flags: append(flags.Flags, flags.S3BucketNameFlag),
			  Action: CreateBucket,
			},
		  },
		},
		{
			Name:    "dynamodb",
			Usage:   "dynamodb operation in s3",
			Subcommands: []cli.Command{
				{
					Name:    "create_metadata_table",
					Aliases: []string{"cmt"},
					Usage:   "create a metadata table",
					Flags: append(flags.Flags, flags.DynamoDBTableNameFlag),
					Action:  CreateMetadataTable,
				},
				{
					Name:    "create_bucket_table",
					Aliases: []string{"cbt"},
					Usage:   "create a bucket table",
					Flags: append(flags.Flags, flags.BucketTableName),
					Action:  CreateMetadataTable,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func CreateBucket(ctx *cli.Context) error {
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

func CreateMetadataTable(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.GetLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	ctx_bg := context.Background()

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	metadataTableName := ctx.String(flags.DynamoDBTableNameFlag.Name)
	if metadataTableName != "" {
		_, err = dynamoClient.CreateTable(ctx_bg, config.AwsClientConfig, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	} else {
		bucketTableName := ctx.String(flags.BucketTableName.Name)
		_, err = dynamoClient.CreateTable(ctx_bg, config.AwsClientConfig, bucketTableName, store.GenerateTableSchema(10, 10, bucketTableName))
	}
	
	if err != nil {
		return err
	}

	return nil
}