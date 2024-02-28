package flags

import (
	"github.com/urfave/cli"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/common/aws"
	"github.com/zero-gravity-labs/zerog-data-avail/common/geth"
	"github.com/zero-gravity-labs/zerog-data-avail/common/logging"
	"github.com/zero-gravity-labs/zerog-data-avail/common/ratelimit"
	"github.com/zero-gravity-labs/zerog-data-avail/common/storage_node"
	server_flags "github.com/zero-gravity-labs/zerog-data-avail/disperser/cmd/apiserver/flags"
	batcher_flags "github.com/zero-gravity-labs/zerog-data-avail/disperser/cmd/batcher/flags"
)

const (
	FlagPrefix   = "combined-server"
	EnvVarPrefix = "COMBINED_SERVER"
)

var (
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "METRICS_HTTP_PORT"),
	}
	EnableMetrics = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: false,
		EnvVar:   common.PrefixEnvVar(EnvVarPrefix, "ENABLE_METRICS"),
	}
	UseMemoryDB = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "use-memory-db"),
		Usage:  "use memory db",
		EnvVar: common.PrefixEnvVar(EnvVarPrefix, "USE_MEMORY_DB"),
	}
)

var RequiredFlags = []cli.Flag{}

var OptionalFlags = []cli.Flag{
	MetricsHTTPPort,
	EnableMetrics,
	UseMemoryDB,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	// combined
	Flags = append(RequiredFlags, OptionalFlags...)
	Flags = append(Flags, logging.CLIFlags(EnvVarPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(EnvVarPrefix)...)
	Flags = append(Flags, aws.ClientFlags(EnvVarPrefix, FlagPrefix)...)
	Flags = append(Flags, storage_node.ClientFlags(EnvVarPrefix, FlagPrefix)...)

	// api server
	Flags = append(Flags, server_flags.RequiredFlags...)
	Flags = append(Flags, server_flags.OptionalFlags...)
	Flags = append(Flags, ratelimit.RatelimiterCLIFlags(server_flags.EnvVarPrefix, server_flags.FlagPrefix)...)

	// batcher
	Flags = append(Flags, batcher_flags.RequiredFlags...)
	Flags = append(Flags, batcher_flags.OptionalFlags...)
}
