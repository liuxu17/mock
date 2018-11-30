package cmd

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/service"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func GenSignedTxDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-signed-tx",
		Short: "generate signed tx data",
		Long: `generate signed tx data
Example:
	mock gen-signed-tx --num {num} --receiver {receiver-address} --faucet-name {faucet-address} --faucet-addr {faucet-addr} --chain-id {chain-id} --node {node-url}
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			num := viper.GetInt(FlagNumSignedTx)
			signedTxData := service.BatchGenSignedTxData(num)

			outputDir := viper.GetString(FlagResOutput)
			err := helper.CreateFolder(outputDir)
			if err != nil {
				panic(err)
			}

			filename := fmt.Sprintf("res_signed_tx_%v", time.Now().Unix())
			filePath := fmt.Sprintf("%v/%v", outputDir, filename)

			err = helper.WriteFile(filePath, []byte(strings.Join(signedTxData, "\n")))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().AddFlagSet(txFlagSet)

	cmd.MarkFlagRequired(FlagNumSignedTx)
	cmd.MarkFlagRequired(FlagReceiverAddr)
	cmd.MarkFlagRequired(FlagFaucetName)
	cmd.MarkFlagRequired(FlagFaucetAddr)
	cmd.MarkPersistentFlagRequired(FlagChainId)
	cmd.MarkPersistentFlagRequired(FlagNodeUrl)

	return cmd
}
