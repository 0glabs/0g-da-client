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
	app.Flags = flags.Flags
	app.Commands = []cli.Command{
		{
			Name:  "bucket",
			Usage: "bucket operation in s3",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "create a new bucket",
					Flags:  append(flags.Flags, flags.S3BucketNameFlag),
					Action: CreateBucket,
				},
				{
					Name:   "delete",
					Usage:  "delete a bucket",
					Flags:  append(flags.Flags, flags.S3BucketNameFlag),
					Action: DeleteBucket,
				},
				{
					Name:   "clear",
					Usage:  "clear a bucket",
					Flags:  append(flags.Flags, flags.S3BucketNameFlag),
					Action: ClearBucket,
				},
			},
		},
		{
			Name:  "dynamodb",
			Usage: "dynamodb operation in s3",
			Subcommands: []cli.Command{
				{
					Name:    "create_metadata_table",
					Aliases: []string{"cmt"},
					Usage:   "create a metadata table",
					Flags:   append(flags.Flags, flags.DynamoDBTableNameFlag),
					Action:  CreateTable,
				},
				{
					Name:    "create_bucket_table",
					Aliases: []string{"cbt"},
					Usage:   "create a bucket table",
					Flags:   append(flags.Flags, flags.BucketTableName),
					Action:  CreateTable,
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "delete a table",
					Flags:   append(flags.Flags, flags.DynamoDBTableNameFlag),
					Action:  DeleteTable,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func CreateBucket(ctx *cli.Context) error {
	config := NewConfig(ctx)

	s3Client, err := getS3Client(config)
	if err != nil {
		return err
	}

	ctx_bg := context.Background()
	bucketName := ctx.String(flags.S3BucketNameFlag.Name)
	log.Println("bucketName: ", bucketName)
	log.Println("region: ", config.AwsClientConfig.Region)
	err = s3Client.CreateBucket(ctx_bg, bucketName, config.AwsClientConfig.Region)

	if err != nil {
		return err
	}

	return nil
}

func DeleteBucket(ctx *cli.Context) error {
	config := NewConfig(ctx)

	s3Client, err := getS3Client(config)
	if err != nil {
		return err
	}

	ctx_bg := context.Background()
	bucketName := ctx.String(flags.S3BucketNameFlag.Name)
	err = s3Client.DeleteBucket(ctx_bg, bucketName)

	if err != nil {
		return err
	}

	return nil
}

func ClearBucket(ctx *cli.Context) error {
	config := NewConfig(ctx)

	s3Client, err := getS3Client(config)
	if err != nil {
		return err
	}

	ctx_bg := context.Background()
	bucketName := ctx.String(flags.S3BucketNameFlag.Name)
	err = s3Client.ClearBucket(ctx_bg, bucketName)

	if err != nil {
		return err
	}

	return nil
}

func CreateTable(ctx *cli.Context) error {
	config := NewConfig(ctx)

	dynamoClient, err := getDynamodbClient(config)
	if err != nil {
		return err
	}

	ctx_bg := context.Background()

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

func DeleteTable(ctx *cli.Context) error {
	config := NewConfig(ctx)

	dynamoClient, err := getDynamodbClient(config)
	ctx_bg := context.Background()

	if err != nil {
		return err
	}

	tableName := ctx.String(flags.DynamoDBTableNameFlag.Name)
	err = dynamoClient.DeleteTable(ctx_bg, tableName)

	if err != nil {
		return err
	}

	return nil
}

func getS3Client(cfg *Config) (*s3.Client, error) {
	logger, err := logging.GetLogger(cfg.LoggerConfig)
	if err != nil {
		return nil, err
	}
	log.Println("cfg.AwsClientConfig: ", cfg.AwsClientConfig)
	s3Client, err := s3.NewClient(cfg.AwsClientConfig, logger)
	if err != nil {
		return nil, err
	}

	return s3Client, nil
}

func getDynamodbClient(cfg *Config) (*dynamodb.Client, error) {
	logger, err := logging.GetLogger(cfg.LoggerConfig)
	if err != nil {
		return nil, err
	}

	dynamoClient, err := dynamodb.NewClient(cfg.AwsClientConfig, logger)
	if err != nil {
		return nil, err
	}

	return dynamoClient, nil
}
