package storage_node

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zgda/common"
)

var (
	StorageNodeURLsFlagName     = "storage.node-urls"
	FlowContractAddressFlagName = "storage.flow-contract"
)

type ClientConfig struct {
	StorageNodeURLs     []string
	FlowContractAddress string
}

func ClientFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:     common.PrefixFlag(flagPrefix, StorageNodeURLsFlagName),
			Usage:    "storage node urls",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "STORAGE_NODE_URLS"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, FlowContractAddressFlagName),
			Usage:    "flow contract address",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "STORAGE_NODE_URLS"),
		},
	}
}

func ReadClientConfig(ctx *cli.Context, flagPrefix string) ClientConfig {
	return ClientConfig{
		StorageNodeURLs:     ctx.GlobalStringSlice(common.PrefixFlag(flagPrefix, StorageNodeURLsFlagName)),
		FlowContractAddress: ctx.GlobalString(common.PrefixFlag(flagPrefix, FlowContractAddressFlagName)),
	}
}
