package cmd

import (
	"bufio"
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/sign"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

func SingleAccSignAndSave() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-signed-tx-separately",
		Short: "generate multi account signed tx data, save the testData",
		Long:  `generate signed tx data`,
		Example: `mock gen-signed-tx-separately \&
--chain-id {chain-id} --home {homeDir} --tps={max broadcast speed} \&
--duration={duration} --bots={num of test node} --account-index={the account index of test nodes}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				subFaucets []conf.SubFaucet
			)
			botsNum := viper.GetInt(FlagBotsNum)
			accountIndex := viper.GetInt(FlagAccountIndex)
			pressureDuration := viper.GetInt(FlagDuration)
			tps := viper.GetInt(FlagTps)

			confHome := viper.GetString(FlagConfDir)
			// create output dir
			outputDir := viper.GetString(FlagResOutput)
			testAccNum := 1 //the special version for stage
			everySignedTxNum := tps * pressureDuration * 60 / botsNum
			log.Printf("Test Prepare sign %d data\n", everySignedTxNum)
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
			if len(subFaucets) <= 0  {
				return fmt.Errorf("can't read sub_faucets config")
			}
			if len(subFaucets) < botsNum {
				return fmt.Errorf("the sub faucet is too few")
			}
			if len(subFaucets) < accountIndex + 1 {
				return fmt.Errorf("accountIndex is out of range")
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
				signedData      conf.SignedData
			)

			log.SetFlags(log.Ldate | log.Lmicroseconds)
			log.Printf("Pressure Test Prepare sign process start\n")

			log.Printf("the no.%d account begin signed data\n", everySignedTxNum)
			accountPriv, err := sign.InitAccountSignProcess(subFaucets[accountIndex].FaucetAddr, subFaucets[accountIndex].Seed)
			signedDataString, err := sign.GenSignTxByTend(everySignedTxNum, accountIndex, conf.ChainId, subFaucets, accountPriv)
			if err != nil {
				return err
			}
			signedData = conf.SignedData{
				SignedDataArray: signedDataString,
			}
			signedDataArray = append(signedDataArray, signedData)
			log.Printf("the no.%d account signed data end\n", accountIndex)

			//set pressure test log to microsecond

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

	cmd.Flags().AddFlagSet(singleSignAndSaveFlagSet)
	//mock gen-signed-tx-separately \&
	//--chain-id {chain-id} --node {node-url} --home {homeDir} --tps={max broadcast speed} \&
	//--duration={duration} --bots={num of test node} --account-index={the account index of test nodes}
	cmd.MarkFlagRequired(FlagAccountIndex)
	cmd.MarkFlagRequired(FlagConfDir)
	cmd.MarkFlagRequired(FlagTps)
	cmd.MarkFlagRequired(FlagBotsNum)
	cmd.MarkFlagRequired(FlagDuration)
	cmd.MarkPersistentFlagRequired(FlagChainId)

	return cmd
}

func BroadCastFromSingleFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast-signed-tx-separately",
		Short: "broadcast signed tx data",
		Long:  `generate signed tx data`,
		Example: `mock broadcast-signed-tx-separately --output {output-dir} --node {node-url} \&
--tps={max broadcast speed} --duration={duration} --bots={num of test node} --commit={block commit time in config}`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// create output dir
			outputDir := viper.GetString(FlagResOutput)

			botsNum := viper.GetInt(FlagBotsNum)
			tps := viper.GetInt(FlagTps)
			commitDuration := viper.GetInt(FlagCommit)

			realDuration := tps * commitDuration / botsNum

			log.SetFlags(log.Ldate | log.Lmicroseconds)
			count := 0
			fin, err := os.OpenFile(outputDir, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("can't find directory in %v\n", outputDir)
			}
			defer fin.Close()

			sc := bufio.NewScanner(fin)
			/*default split the file use '\n'*/
			log.Printf("Pressure Test Start !!!!!!!!!!!\n")
			timeTemp := time.Now()
			log.Printf("test account start txs\n")
			for sc.Scan() {
				count++
				_, err = sign.BroadcastTx(sc.Text())
				if count % realDuration == 0 {
					if timeTemp.Add(time.Second * time.Duration(commitDuration)).After(time.Now()) {
						time.Sleep(timeTemp.Add(time.Second * time.Duration(commitDuration)).Sub(time.Now()))
						log.Printf("test broadcast %d\n", count)
					}
					timeTemp = time.Now()
				}
			}
			if err := sc.Err(); err != nil{
				fmt.Println("An error has happened, when we run buf scanner")
				return err
			}
			log.Printf("%v: all test is over\n", "PressureTest")
			return nil
		},
	}

	cmd.Flags().AddFlagSet(broadcastFlagSet)
	//mock gen-signed-tx-separately \&
	//--chain-id {chain-id} --node {node-url} --home {homeDir} --tps={max broadcast speed} \&
	//--duration={duration} --bots={num of test node} --account-index={the account index of test nodes}
	cmd.MarkFlagRequired(FlagTps)
	cmd.MarkFlagRequired(FlagBotsNum)
	cmd.MarkFlagRequired(FlagResOutput)
	cmd.MarkFlagRequired(FlagCommit)
	cmd.MarkPersistentFlagRequired(FlagNodeUrl)

	return cmd
}
