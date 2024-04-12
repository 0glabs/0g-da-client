package geth

import (
	"time"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/urfave/cli"
)

var (
	rpcUrlFlagName                 = "chain.rpc"
	cosmosGrpcUrlFlagName          = "chain.cosmos-grpc"
	privateKeyFlagName             = "chain.private-key"
	numConfirmationsFlagName       = "chain.num-confirmations"
	txGasLimitFlagName             = "chain.gas-limit"
	receiptPollingRoundsFlagName   = "chain.receipt-wait-rounds"
	receiptPollingIntervalFlagName = "chain.receipt-wait-interval"
)

type EthClientConfig struct {
	RPCURL                 string
	CosmosGrpc             string
	PrivateKeyString       string
	NumConfirmations       int
	TxGasLimit             int
	ReceiptPollingRounds   uint
	ReceiptPollingInterval time.Duration
}

func EthClientFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     cosmosGrpcUrlFlagName,
			Usage:    "Cosmos grpc ",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "COSMOS_GRPC"),
		},
		cli.StringFlag{
			Name:     rpcUrlFlagName,
			Usage:    "Chain rpc",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "CHAIN_RPC"),
		},
		cli.StringFlag{
			Name:     privateKeyFlagName,
			Usage:    "Ethereum private key for disperser",
			Required: false,
			Value:    "0000000000000000000000000000000000000000000000000000000000000000",
			EnvVar:   common.PrefixEnvVar(envPrefix, "PRIVATE_KEY"),
		},
		cli.IntFlag{
			Name:     numConfirmationsFlagName,
			Usage:    "Number of confirmations to wait for",
			Required: false,
			Value:    0,
			EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_CONFIRMATIONS"),
		},
		cli.IntFlag{
			Name:     txGasLimitFlagName,
			Usage:    "Gas limit for transaction",
			Required: false,
			Value:    0,
			EnvVar:   common.PrefixEnvVar(envPrefix, "TX_GAS_LIMIT"),
		},
		cli.UintFlag{
			Name:     receiptPollingRoundsFlagName,
			Usage:    "Rounds of receipt polling",
			Required: false,
			Value:    60,
			EnvVar:   common.PrefixEnvVar(envPrefix, "RECEIPT_POLLING_ROUNDS"),
		},
		cli.DurationFlag{
			Name:     receiptPollingIntervalFlagName,
			Usage:    "Interval of receipt polling",
			Required: false,
			Value:    time.Second,
			EnvVar:   common.PrefixEnvVar(envPrefix, "RECEIPT_POLLING_INTERVAL"),
		},
	}
}

func ReadEthClientConfig(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	cfg.CosmosGrpc = ctx.GlobalString(cosmosGrpcUrlFlagName)
	cfg.PrivateKeyString = ctx.GlobalString(privateKeyFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	cfg.TxGasLimit = ctx.GlobalInt(txGasLimitFlagName)
	cfg.ReceiptPollingRounds = ctx.GlobalUint(receiptPollingRoundsFlagName)
	cfg.ReceiptPollingInterval = ctx.GlobalDuration(receiptPollingIntervalFlagName)
	return cfg
}
