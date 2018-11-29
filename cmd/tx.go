package cmd

import (
	"github.com/spf13/cobra"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/spf13/viper"
	"github.com/kaifei-bianjie/mock/service"
	"fmt"
	"time"
	"github.com/kaifei-bianjie/mock/util/helper"
	"strings"
	"os"
)

func GenSignedTxDataCmd() *cobra.Command  {
	cmd := &cobra.Command{
		Use:   "gen-signed-tx",
		Short: "generate signed tx data",
		Long:  `generate signed tx data`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			num := viper.GetInt(FlagNumSignedTx)
			signedTxData := service.BatchGenSignedTxData(num)

			homeDir := viper.GetString(FlagResSignedTxFileHome)
			filename := fmt.Sprintf("res_signed_tx_%v.txt", time.Now().Unix())
			filePath := fmt.Sprintf("%v/%v", homeDir, filename)

			err := helper.WriteFile(filePath, []byte(strings.Join(signedTxData, "\n")))
			if err != nil {
				return err
			}

			return nil
		},
	}


	cmd.LocalFlags().StringP(FlagResSignedTxFileHome, "", os.ExpandEnv("$HOME"), "directory of result file which content signed tx data")
	cmd.LocalFlags().IntP(FlagNumSignedTx, "", 1, "num of signed tx which need to generated")
	cmd.LocalFlags().IntVarP(&conf.BlockInterval, FlagBlockInterval, "", 5, "block interval")
	cmd.LocalFlags().StringVarP(&conf.DefaultReceiverAddr, FlagReceiverAddr, "", "", "receiver address")

	cmd.MarkFlagRequired(FlagReceiverAddr)

	return cmd
}