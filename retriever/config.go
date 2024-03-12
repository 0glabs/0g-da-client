package retriever

import (
	"github.com/0glabs/0g-data-avail/common/logging"
	"github.com/0glabs/0g-data-avail/common/storage_node"
	"github.com/0glabs/0g-data-avail/core/encoding"
	"github.com/0glabs/0g-data-avail/retriever/flags"
	"github.com/urfave/cli"
)

type Config struct {
	EncoderConfig     encoding.EncoderConfig
	LoggerConfig      logging.Config
	StorageNodeConfig storage_node.ClientConfig
	MetricsConfig     MetricsConfig

	NumConnections int
}

func NewConfig(ctx *cli.Context) *Config {
	return &Config{
		EncoderConfig: encoding.ReadCLIConfig(ctx),
		LoggerConfig:  logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		MetricsConfig: MetricsConfig{
			HTTPPort: ctx.GlobalString(flags.MetricsHTTPPortFlag.Name),
		},
		StorageNodeConfig: storage_node.ReadClientConfig(ctx, flags.FlagPrefix),
		NumConnections:    ctx.Int(flags.NumConnectionsFlag.Name),
	}
}
