package main

import (
	"github.com/spf13/cobra"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mock",
		Short: "mock test",
	}
)

func main() {

	rootCmd.AddCommand(
		cmd.FaucetInitCmd(),
		cmd.GenSignedTxDataCmd(),
	)

	rootCmd.PersistentFlags().StringVarP(&conf.NodeUrl, cmd.FlagNodeUrl, "", "", "http://localhost:1317")
	rootCmd.PersistentFlags().StringVarP(&conf.ChainId, cmd.FlagChainId, "", "", "testnet")

	rootCmd.MarkFlagRequired(cmd.FlagChainId)
	rootCmd.MarkFlagRequired(cmd.FlagNodeUrl)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
