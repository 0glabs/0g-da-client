package storage_node

import (
	"github.com/0glabs/0g-da-client/common"
	"github.com/urfave/cli"
)

var (
	KVDBPathFlagName     = "storage.kv-db-path"
	TimeToExpireFlagName = "storage.time-to-expire"
)

type ClientConfig struct {
	KvDbPath     string
	TimeToExpire uint
}

func ClientFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
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
			Value:    5184000, // 60 days
			EnvVar:   common.PrefixEnvVar(envPrefix, "TimeToExpire"),
		},
	}
}

func ReadClientConfig(ctx *cli.Context, flagPrefix string) ClientConfig {
	return ClientConfig{
		KvDbPath:     ctx.GlobalString(common.PrefixFlag(flagPrefix, KVDBPathFlagName)),
		TimeToExpire: ctx.GlobalUint(common.PrefixFlag(flagPrefix, TimeToExpireFlagName)),
	}
}
