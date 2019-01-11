package cmd

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/service"
	"github.com/kaifei-bianjie/mock/sign"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

func MultiAccGenSignedTxDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-multi-signed-tx",
		Short: "generate multi account signed tx data",
		Long:  `batch generate signed tx data`,
		Example: `mock gen-multi-signed-tx --num {num} --faucet-name {faucet-name}\&
--chain-id {chain-id} --node {node-url} --home {homeDir} --output {optional}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				subFaucets []conf.SubFaucet
			)
			confHome := viper.GetString(FlagConfDir)
			// create output dir
			outputDir := viper.GetString(FlagResOutput)
			testAccNum := viper.GetInt(FlagTestAccountNum)
			everySignedTxNum := viper.GetInt(FlagEveryNumSignedTx)
			// read config from config file
			viper.SetConfigName("config")
			viper.AddConfigPath(confHome)
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}
			err = viper.UnmarshalKey("sub_faucets", &subFaucets)
			if err != nil {
				return err
			}
			if len(subFaucets) <= 0 {
				return fmt.Errorf("can't read sub_faucets config")
			}

			err = helper.CreateFolder(outputDir)
			if err != nil {
				return err
			}

			// read config from config file
			// check config file if exists
			exists, err := helper.CheckFileExist(confHome)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("can't find keys directory in %v\n", confHome)
			}

			var (
				signedDataArray []conf.SignedData
				signedData 		conf.SignedData
			)
			for i := 0; i < testAccNum; i++ {
				signedData = conf.SignedData{
					SignedDataArray: service.MultiBatchGenSignedTxData(everySignedTxNum, subFaucets[i].FaucetName, subFaucets[i].FaucetAddr, confHome, subFaucets),
				}
				signedDataArray = append(signedDataArray, signedData)
				log.Printf("the no.%d account signed data end\n",i)
			}

			//set pressure test log to microsecond
			log.SetFlags(log.Ldate | log.Lmicroseconds)
			log.Printf("Pressure Test Start !!!!!!!!!!!\n")
			counterChan := make(chan types.TestPressData, 100000)
			for i := 0; i < testAccNum; i++ {
				go sign.BroadcastTxForAccount(signedDataArray, i, counterChan)
			}

			counter := 0
			for {
				testInfo := <-counterChan
				//log.Printf("test account no.%d succeed %d txs, sum is %d\n", testInfo.AccountIndex, testInfo.SuccessIndex, everySignedTxNum)
				counter++
				if counter == testAccNum {
					log.Printf("%v: all test is over %d\n", "PressureTest", testInfo.AccountIndex)
					break
				}
			}

			filename := fmt.Sprintf("res_signed_tx_%v", time.Now().Unix())
			filePath := fmt.Sprintf("%v/%v", outputDir, filename)
			var signedTxData []string
			if len(signedDataArray) > 0 {
				// write result to file
				for i := 0; i < everySignedTxNum; i++ {
					for j := 0; j < testAccNum; j++ {
						signedTxData = append(signedTxData, signedDataArray[j].SignedDataArray[i])
					}
				}
			} else {
				fmt.Println("no signed tx data")
			}

			err = helper.WriteFile(filePath, []byte(strings.Join(signedTxData, "\n")))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().AddFlagSet(multiTxFlagSet)
	//--num {num} --faucet-name {faucet-name}\&
	//--chain-id {chain-id} --node {node-url} --home {homeDir} --output {optional}
	cmd.MarkFlagRequired(FlagEveryNumSignedTx)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagTestAccountNum)

	cmd.MarkPersistentFlagRequired(FlagChainId)
	cmd.MarkPersistentFlagRequired(FlagNodeUrl)

	return cmd
}



func MultiAccSignDirectly() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-multi-signed-tx-directly",
		Short: "generate multi account signed tx data directly",
		Long:  `batch generate signed tx data`,
		Example: `mock gen-multi-signed-tx-directly --num {num} --faucet-name {faucet-name}\&
--chain-id {chain-id} --node {node-url} --home {homeDir} --output {optional}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				subFaucets []conf.SubFaucet
			)
			confHome := viper.GetString(FlagConfDir)
			// create output dir
			outputDir := viper.GetString(FlagResOutput)
			testAccNum := viper.GetInt(FlagTestAccountNum)
			everySignedTxNum := viper.GetInt(FlagEveryNumSignedTx)
			// read config from config file
			viper.SetConfigName("config")
			viper.AddConfigPath(confHome)
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}
			err = viper.UnmarshalKey("sub_faucets", &subFaucets)
			if err != nil {
				return err
			}
			if len(subFaucets) <= 0 {
				return fmt.Errorf("can't read sub_faucets config")
			}
			err = helper.CreateFolder(outputDir)
			if err != nil {
				return err
			}
			// read config from config file
			// check config file if exists
			exists, err := helper.CheckFileExist(confHome)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("can't find keys directory in %v\n", confHome)
			}

			var (
				signedDataArray []conf.SignedData
				signedData 		conf.SignedData
			)


			for i := 0; i < testAccNum; i++ {
				log.Printf("the no.%d account begin signed data\n",i)
				accountPriv, err := sign.InitAccountSignProcess(subFaucets[i].FaucetAddr, subFaucets[i].Seed)
				signedDataString, err := sign.GenSignTxByTend(everySignedTxNum, conf.ChainId, subFaucets, accountPriv)
				if err != nil {
					return err
				}
				signedData = conf.SignedData{
					SignedDataArray: signedDataString,
				}
				signedDataArray = append(signedDataArray, signedData)
				log.Printf("the no.%d account signed data end\n",i)
			}

			//set pressure test log to microsecond
			log.SetFlags(log.Ldate | log.Lmicroseconds)
			log.Printf("Pressure Test Start !!!!!!!!!!!\n")
			counterChan := make(chan types.TestPressData, 100000)
			for i := 0; i < testAccNum; i++ {
				go sign.BroadcastTxForAccount(signedDataArray, i, counterChan)
			}

			counter := 0
			for {
				testInfo := <-counterChan
				//log.Printf("test account no.%d succeed %d txs, sum is %d\n", testInfo.AccountIndex, testInfo.SuccessIndex, everySignedTxNum)
				counter++
				if counter == testAccNum {
					log.Printf("%v: all test is over %d\n", "PressureTest", testInfo.AccountIndex)
					break
				}
			}

			filename := fmt.Sprintf("res_signed_tx_%v", time.Now().Unix())
			filePath := fmt.Sprintf("%v/%v", outputDir, filename)
			var signedTxData []string
			if len(signedDataArray) > 0 {
				// write result to file
				for i := 0; i < everySignedTxNum; i++ {
					for j := 0; j < testAccNum; j++ {
						signedTxData = append(signedTxData, signedDataArray[j].SignedDataArray[i])
					}
				}
			} else {
				fmt.Println("no signed tx data")
			}

			err = helper.WriteFile(filePath, []byte(strings.Join(signedTxData, "\n")))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().AddFlagSet(multiTxDirectlyFlagSet)
	//--num {num} --faucet-name {faucet-name}\&
	//--chain-id {chain-id} --node {node-url} --home {homeDir} --output {optional}
	cmd.MarkFlagRequired(FlagEveryNumSignedTx)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagTestAccountNum)

	cmd.MarkPersistentFlagRequired(FlagChainId)
	cmd.MarkPersistentFlagRequired(FlagNodeUrl)

	return cmd
}
