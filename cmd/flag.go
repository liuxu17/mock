package cmd

import (
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/spf13/pflag"
	"os"
)

const (
	FlagNodeUrl = "node"
	FlagChainId = "chain-id"
	FlagConfDir = "home"

	FlagFaucetSeed      = "seed"
	FlagFaucetName      = "faucet-name"
	FlagSubFaucetAccNum = "sub-faucet-num"

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
	rootFlagSet.String(FlagConfDir, conf.DefaultHome, "directory of config file")

	faucetFlagSet.StringVarP(&conf.FaucetName, FlagFaucetName, "", "", "faucet name")
	faucetFlagSet.StringVarP(&conf.FaucetSeed, FlagFaucetSeed, "", "", "seed")
	faucetFlagSet.String(FlagConfDir, conf.DefaultHome, "directory for save config data")
	faucetFlagSet.Int(FlagSubFaucetAccNum, 10, "num of sub faucet want to create, shouldn't greater than 10")

	txFlagSet.Int(FlagNumSignedTx, 0, "num of signed tx which need to generated")
	txFlagSet.StringVar(&conf.DefaultReceiverAddr, FlagReceiverAddr, "", "receiver address")
	txFlagSet.IntVar(&conf.BlockInterval, FlagBlockInterval, 5, "block interval")
	txFlagSet.String(FlagResOutput, os.ExpandEnv("$HOME")+"/output", "output directory of result file which content signed tx data")
	txFlagSet.String(FlagConfDir, conf.DefaultHome, "directory of config file")
}
