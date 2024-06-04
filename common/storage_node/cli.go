package storage_node

import (
	"github.com/0glabs/0g-data-avail/common"
	"github.com/urfave/cli"
)

var (
	DAEntranceContractAddressFlagName = "storage.da-entrance-contract"
	DASignersContractAddressFlagName  = "storage.da-signers-contract"
	KVDBPathFlagName                  = "storage.kv-db-path"
	TimeToExpireFlagName              = "storage.time-to-expire"
	UploadTaskSizeFlagName            = "storage.upload-task-size"
)

type ClientConfig struct {
	DAEntranceContractAddress string
	DASignersContractAddress  string
	KvDbPath                  string
	TimeToExpire              uint
	UploadTaskSize            uint
}

func ClientFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, DAEntranceContractAddressFlagName),
			Usage:    "DAEntrance contract address",
			Required: false,
			Value:    "0x0000000000000000000000000000000000000000",
			EnvVar:   common.PrefixEnvVar(envPrefix, "DAENTRANCE_CONTRACT_ADDRESS"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, DASignersContractAddressFlagName),
			Usage:    "DASigners contract address",
			Required: false,
			Value:    "0x0000000000000000000000000000000000000000",
			EnvVar:   common.PrefixEnvVar(envPrefix, "DASIGNERS_CONTRACT_ADDRESS"),
		},
		cli.StringFlag{
			Name:     common.PrefixFlag(flagPrefix, KVDBPathFlagName),
			Usage:    "kv db path",
			Required: false,
			Value:    "",
			EnvVar:   common.PrefixEnvVar(envPrefix, "KV_DB_PATH"),
		},
		cli.UintFlag{
			Name:     common.PrefixFlag(flagPrefix, TimeToExpireFlagName),
			Usage:    "time to expire",
			Required: false,
			Value:    5184000,
			EnvVar:   common.PrefixEnvVar(envPrefix, "TimeToExpire"),
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
	return ClientConfig{
		DAEntranceContractAddress: ctx.GlobalString(common.PrefixFlag(flagPrefix, DAEntranceContractAddressFlagName)),
		DASignersContractAddress:  ctx.GlobalString(common.PrefixFlag(flagPrefix, DASignersContractAddressFlagName)),
		UploadTaskSize:            ctx.GlobalUint(common.PrefixFlag(flagPrefix, UploadTaskSizeFlagName)),
		KvDbPath:                  ctx.GlobalString(common.PrefixFlag(flagPrefix, KVDBPathFlagName)),
		TimeToExpire:              ctx.GlobalUint(common.PrefixFlag(flagPrefix, TimeToExpireFlagName)),
	}
}
