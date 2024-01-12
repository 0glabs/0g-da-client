package aws

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
)

var (
	RegionFlagName          = "aws.region"
	AccessKeyIdFlagName     = "aws.access-key-id"
	SecretAccessKeyFlagName = "aws.secret-access-key"
	EndpointURLFlagName     = "aws.endpoint-url"
)

type ClientConfig struct {
	Region          string
	AccessKey       string
	SecretAccessKey string
	EndpointURL     string
	UseMemDb        bool
}

func ClientFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, RegionFlagName),
			Usage:    "AWS Region",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_REGION"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, AccessKeyIdFlagName),
			Usage:    "AWS Access Key Id",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_ACCESS_KEY_ID"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, SecretAccessKeyFlagName),
			Usage:    "AWS Secret Access Key",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_SECRET_ACCESS_KEY"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, EndpointURLFlagName),
			Usage:    "AWS Endpoint URL",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "AWS_ENDPOINT_URL"),
		},
	}
}

func ReadClientConfig(ctx *cli.Context, flagPrefix string) ClientConfig {
	return ClientConfig{
		Region:          ctx.GlobalString(common.PrefixFlag(flagPrefix, RegionFlagName)),
		AccessKey:       ctx.GlobalString(common.PrefixFlag(flagPrefix, AccessKeyIdFlagName)),
		SecretAccessKey: ctx.GlobalString(common.PrefixFlag(flagPrefix, SecretAccessKeyFlagName)),
		EndpointURL:     ctx.GlobalString(common.PrefixFlag(flagPrefix, EndpointURLFlagName)),
	}
}
