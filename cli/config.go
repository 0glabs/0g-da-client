package main

import (
	"github.com/0glabs/0g-da-client/cli/flags"
	"github.com/0glabs/0g-da-client/common/aws"
	"github.com/0glabs/0g-da-client/common/logging"
	"github.com/urfave/cli"
)

type Config struct {
	AwsClientConfig aws.ClientConfig
	LoggerConfig    logging.Config
}

func NewConfig(ctx *cli.Context) *Config {
	return &Config{
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		LoggerConfig:    logging.ReadCLIConfig(ctx, flags.FlagPrefix),
	}
}
