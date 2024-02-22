package main

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
	"github.com/zero-gravity-labs/zerog-data-avail/cli/flags"
)

type Config struct {
	AwsClientConfig   aws.ClientConfig
	LoggerConfig      logging.Config
}

func NewConfig(ctx *cli.Context) (Config, error) {
	config := Config{
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		LoggerConfig: logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		
	}
	return config, nil
}
