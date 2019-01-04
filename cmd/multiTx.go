package cmd

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/service"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func MultiAccGenSignedTxDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-multi-signed-tx",
		Short: "generate multi account signed tx data",
		Long:  `batch generate signed tx data`,
		Example: `mock gen-multi-signed-tx --num {num} --faucet-name {faucet-name}\&
--chain-id {chain-id} --node {node-url} --home {homeDir} --output {optional} --address {faucetAddress}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// create output dir
			outputDir := viper.GetString(FlagResOutput)
			faucetAddr := viper.GetString(FlagFaucetAddress)
			err := helper.CreateFolder(outputDir)
			if err != nil {
				return err
			}

			// read config from config file
			// check config file if exists
			confHome := viper.GetString(FlagConfDir)

			exists, err := helper.CheckFileExist(confHome)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("can't find keys directory in %v\n", confHome)
			}

			// generate signed tx
			num := viper.GetInt(FlagNumSignedTx)
			if num > 100000 {
				return fmt.Errorf("num should less than 100000\n")
			}

			name := viper.GetString(FlagFaucetName)

			signedTxData := service.SingleBatchGenSignedTxData(num, name, constants.KeyPassword, confHome, faucetAddr)

			if len(signedTxData) > 0 {
				// write result to file
				filename := fmt.Sprintf("res_signed_tx_%v", time.Now().Unix())
				filePath := fmt.Sprintf("%v/%v", outputDir, filename)
				err = helper.WriteFile(filePath, []byte(strings.Join(signedTxData, "\n")))
				if err != nil {
					return err
				}
			} else {
				fmt.Println("no signed tx data")
			}
			return nil
		},
	}

	cmd.Flags().AddFlagSet(multiTxFlagSet)
	//--num {num} --faucet-name {faucet-name}\&
	//--chain-id {chain-id} --node {node-url} --home {homeDir} --output {optional}
	cmd.MarkFlagRequired(FlagNumSignedTx)
	cmd.MarkFlagRequired(FlagFaucetName)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagFaucetName)
	cmd.MarkPersistentFlagRequired(FlagChainId)
	cmd.MarkPersistentFlagRequired(FlagNodeUrl)

	return cmd
}
