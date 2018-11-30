package cmd

import (
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/spf13/pflag"
	"os"
)

const (
	FlagNodeUrl = "node"
	FlagChainId = "chain-id"
	FlagHome    = "home"

	FlagFaucetSeed = "seed"
	FlagFaucetName = "faucet-name"
	FlagFaucetAddr = "faucet-addr"

	FlagBlockInterval = "block-interval"
	FlagReceiverAddr  = "receiver"
	FlagNumSignedTx   = "num"
	FlagResOutput     = "output"
)

var (
	rootFlagSet   = pflag.NewFlagSet("", pflag.ContinueOnError)
	faucetFlagSet = pflag.NewFlagSet("", pflag.ContinueOnError)
	txFlagSet     = pflag.NewFlagSet("", pflag.ContinueOnError)
)

func init() {
	rootFlagSet.StringVar(&conf.ChainId, FlagChainId, "", "chain id")
	rootFlagSet.StringVar(&conf.NodeUrl, FlagNodeUrl, "http://localhost:1317", "lcd url")
	rootFlagSet.String(FlagHome, conf.DefaultHome, "directory of config file")

	faucetFlagSet.StringVarP(&conf.FaucetName, FlagFaucetName,"", "", "faucet name")
	faucetFlagSet.StringVarP(&conf.FaucetSeed, FlagFaucetSeed, "", "", "seed")

	txFlagSet.String(FlagResOutput, os.ExpandEnv("$HOME")+"/output", "output directory of result file which content signed tx data")
	txFlagSet.Int(FlagNumSignedTx, 0, "num of signed tx which need to generated")
	txFlagSet.IntVar(&conf.BlockInterval, FlagBlockInterval, 5, "block interval")
	txFlagSet.StringVar(&conf.DefaultReceiverAddr, FlagReceiverAddr, "", "receiver address")
	txFlagSet.StringVar(&conf.FaucetName, FlagFaucetName, "", "faucet name")
	txFlagSet.StringVar(&conf.FaucetAddress, FlagFaucetAddr, "", "faucet address")
}
