package storage_node

import (
	"github.com/0glabs/0g-data-avail/common"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var (
	StorageNodeURLsFlagName     = "storage.node-url"
	KVNodeURLFlagName           = "storage.kv-url"
	KVStreamIDFlagName          = "storage.kv-stream-id"
	FlowContractAddressFlagName = "storage.flow-contract"
	UploadTaskSizeFlagName      = "storage.upload-task-size"
)

type ClientConfig struct {
	StorageNodeURLs     []string
	FlowContractAddress string
	KVNodeURL           string
	KVStreamId          eth_common.Hash
	UploadTaskSize      uint
}

func ClientFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:     common.PrefixFlag(flagPrefix, StorageNodeURLsFlagName),
			Usage:    "storage node urls",
			Required: false,
			Value:    nil,
			EnvVar:   common.PrefixEnvVar(envPrefix, "STORAGE_NODE_URLS"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, FlowContractAddressFlagName),
			Usage:    "flow contract address",
			Required: false,
			Value:    "0x0000000000000000000000000000000000000000",
			EnvVar:   common.PrefixEnvVar(envPrefix, "STORAGE_NODE_URLS"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, KVNodeURLFlagName),
			Usage:    "kv node url",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "KV_NODE_URL"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, KVStreamIDFlagName),
			Usage:    "kv stream id",
			Required: false,
			Value:    "0000000000000000000000000000000000000000000000000000000000000000",
			EnvVar:   common.PrefixEnvVar(envPrefix, "KV_NODE_URL"),
		},
		cli.UintFlag{
			Name:     common.PrefixFlag(flagPrefix, UploadTaskSizeFlagName),
			Usage:    "number of segments in single upload rpc request",
			Required: false,
			Value:    10,
			EnvVar:   common.PrefixEnvVar(envPrefix, "UPLOAD_TASK_SIZE"),
		},
	}
}

func ReadClientConfig(ctx *cli.Context, flagPrefix string) ClientConfig {
	streamId := eth_common.HexToHash(ctx.GlobalString(common.PrefixFlag(flagPrefix, KVStreamIDFlagName)))
	return ClientConfig{
		StorageNodeURLs:     ctx.GlobalStringSlice(common.PrefixFlag(flagPrefix, StorageNodeURLsFlagName)),
		FlowContractAddress: ctx.GlobalString(common.PrefixFlag(flagPrefix, FlowContractAddressFlagName)),
		KVNodeURL:           ctx.GlobalString(common.PrefixFlag(flagPrefix, KVNodeURLFlagName)),
		KVStreamId:          streamId,
		UploadTaskSize:      ctx.GlobalUint(common.PrefixFlag(flagPrefix, UploadTaskSizeFlagName)),
	}
}
