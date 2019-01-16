package key

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/sign"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"log"
	"strconv"
	"strings"
	"time"
)

// create faucet sub accounts and transfer token to them.
// faucetName and faucetPasswd specified by flag,
// subAccNum is the num of sub account which need create.
func CreateFaucetSubAccount(faucetName, faucetPasswd, faucetAddr string, subAccNum int) ([]types.AccountInfo, error) {
	var (
		method               = "CreateFaucetSubAccount"
		createdAccs, subAccs []types.AccountInfo
	)

	keyChan := make(chan types.AccountInfo)

	// create sub account
	for i := 1; i <= subAccNum; i++ {
		keyName := fmt.Sprintf("%v_%v", faucetName, i)
		go CreateKey(keyName, keyChan)
	}

	counter := 0
	for {
		accInfo := <-keyChan
		if accInfo.Address != "" {
			createdAccs = append(createdAccs, accInfo)
		}
		counter++
		if counter == subAccNum {
			log.Printf("%v: all create sub faucet key goroutine over\n", method)
			log.Printf("%v: except create %v accounts, successful create %v accounts",
				method, subAccNum, len(createdAccs))
			break
		}
	}

	// distribute token

	// get sender info
	senderInfo := types.AccountInfo{
		LocalAccountName: faucetName,
		Password:         faucetPasswd,
		Address:          faucetAddr,
	}
	accInfo, err := account.GetAccountInfo(senderInfo.Address)
	if err != nil {
		log.Printf("%v: get faucet info fail: %v\n", method, err)
		return subAccs, err
	}
	senderInfo.AccountNumber = accInfo.AccountNumber
	senderSequence, err := helper.ConvertStrToInt64(accInfo.Sequence)
	if err != nil {
		log.Printf("%v: convert sequence to int64 fail: %v\n", method, err)
		return subAccs, err
	}

	// get transfer amount which equal senderBalance / subAccNum
	amt, err := parseCoins(accInfo.Coins)
	if err != nil {
		log.Printf("%v: parse coin failed: %v\n", method, err)
		return subAccs, err
	}
	transferAmt := fmt.Sprintf("%v%s", parseFloat64ToStr(amt/float64(subAccNum+1)), constants.Denom)

	// distribute token to created accounts
	for _, acc := range createdAccs {
		senderInfo.Sequence = fmt.Sprintf("%v", senderSequence)
		acc, err := DistributeToken(senderInfo, acc, transferAmt)
		if err != nil {
			log.Printf("%v: distribute token to %v failed: %v\n",
				method, acc.LocalAccountName, err)
		} else {
			subAccs = append(subAccs, acc)
			senderSequence += 1
		}
	}

	return subAccs, err
}

func CreateSubAccount(faucetName, faucetPasswd, confHome string, subAccNum int, faucetAddr string) ([]types.AccountInfo, error) {
	var (
		method               = "CreateFaucetSubAccount"
		createdAccs, subAccs []types.AccountInfo
	)

	keyChan := make(chan types.AccountInfo, 100000)

	// create sub account
	for i := 1; i <= subAccNum; i++ {
		keyName := fmt.Sprintf("%v_%v", faucetName, i)
		CreateKeyByCli(keyName, keyChan, confHome)
	}

	counter := 0
	for {
		accInfo := <-keyChan
		if accInfo.Address != "" {
			createdAccs = append(createdAccs, accInfo)
		}
		counter++
		if counter == subAccNum {
			log.Printf("%v: all create sub faucet key goroutine over\n", method)
			log.Printf("%v: except create %v accounts, successful create %v accounts",
				method, subAccNum, len(createdAccs))
			break
		}
	}

	// distribute token

	// get sender info
	senderInfo := types.AccountInfo{
		LocalAccountName: faucetName,
		Password:         faucetPasswd,
		Address:          faucetAddr,
	}
	accInfo, err := account.GetAccountInfo(senderInfo.Address)
	if err != nil {
		log.Printf("%v: get faucet info fail: %v\n", method, err)
		return subAccs, err
	}
	senderInfo.AccountNumber = accInfo.AccountNumber
	senderSequence, err := helper.ConvertStrToInt64(accInfo.Sequence)
	if err != nil {
		log.Printf("%v: convert sequence to int64 fail: %v\n", method, err)
		return subAccs, err
	}

	// get transfer amount which equal senderBalance / subAccNum
	amt, err := parseCoins(accInfo.Coins)
	if err != nil {
		log.Printf("%v: parse coin failed: %v\n", method, err)
		return subAccs, err
	}
	transferAmt := fmt.Sprintf("%v%s", parseFloat64ToStr(amt/float64(subAccNum+1)), constants.Denom)

	// distribute token to created accounts
	for _, acc := range createdAccs {
		senderInfo.Sequence = fmt.Sprintf("%v", senderSequence)
		acc, err := DistributeToken(senderInfo, acc, transferAmt)
		if err != nil {
			log.Printf("%v: distribute token to %v failed: %v\n",
				method, acc.LocalAccountName, err)
		} else {
			subAccs = append(subAccs, acc)
			senderSequence += 1
		}
	}

	return subAccs, err
}

func CreateSubAccountByBroadCast(faucetName, faucetPasswd, confHome string, subAccNum int, faucetAddr string) ([]types.AccountInfo, error) {
	var (
		method  = "CreateFaucetSubAccount"
		subAccs []types.AccountInfo
	)
	resChan := make(chan types.GenSignedTxDataRes, 100000)

	// get sender info
	senderInfo := types.AccountInfo{
		LocalAccountName: faucetName,
		Password:         faucetPasswd,
		Address:          faucetAddr,
	}
	accInfo, err := account.GetAccountInfo(senderInfo.Address)
	if err != nil {
		log.Printf("%v: get faucet info fail: %v\n", method, err)
		return subAccs, err
	}
	senderInfo.AccountNumber = accInfo.AccountNumber
	if err != nil {
		log.Printf("%v: convert sequence to int64 fail: %v\n", method, err)
		return subAccs, err
	}

	// get transfer amount which equal senderBalance / subAccNum
	amt, err := parseCoins(accInfo.Coins)
	if err != nil {
		log.Printf("%v: parse coin failed: %v\n", method, err)
		return subAccs, err
	}
	transferAmt := fmt.Sprintf("%v%s", parseFloat64ToStr(amt/float64(subAccNum+1)), constants.Denom)

	// distribute token to created accounts

	senderInfos, err := NewAccountSingle(subAccNum, confHome, faucetName, faucetAddr)
	if err != nil {
		// TODO: handle err
	}

	lens := len(senderInfos)
	var (
		signedTxData [100000]string
	)

	if lens > 0 {
		log.Printf("%v: now use %v goroutine to gen signed data\n",
			method, lens)

		for i, senderInfoEx := range senderInfos {
			go sign.GenSignedTxDataByAmountAndFaucet(transferAmt, faucetAddr, senderInfoEx, conf.DefaultReceiverAddr, resChan, i)
			if (i % 10) == 0 {
				time.Sleep(time.Second * constants.SignTxDelaySec)
			}
		}

		counter := 0
		for {
			res := <-resChan
			counter++
			if res.Res != "" {
				log.Printf("%v: successed, goroutine %v gen signed tx data. now left %v goroutine\n",
					method, res.ChanNum, lens-counter)
				signedTxData[res.ChanNum] = res.Res
				//signedTxData = append(signedTxData, res.Res)
			} else {
				log.Printf("%v: failed, goroutine %v gen signed tx data. now left %v goroutine\n",
					method, res.ChanNum, lens-counter)
			}

			if counter == lens {
				log.Printf("%v: all sign tx goroutine over\n", method)
				break
			}
		}
	} else {
		log.Printf("%v: no signed tx data\n", method)
	}

	var signedTxDataReturn []string
	for j := 0; j < lens; j++ {
		signedTxDataReturn = append(signedTxDataReturn, signedTxData[j])
	}

	var counter int
	for i := 0; i < len(signedTxDataReturn); i++ {
		_, err = sign.BroadcastTx(signedTxDataReturn[i])
		// handle response
		if err != nil {
			log.Printf("test faucet account succeed %d txs, sum is %d\n", counter, subAccNum)
			break
		} else {
			counter = counter + 1
		}
	}
	log.Printf("test faucet account succeed %d txs, sum is %d\n", counter, subAccNum)

	for _, acc := range senderInfos {
		acc = types.AccountInfo{
			Password:         acc.Password,
			LocalAccountName: acc.AccountName,
			Address:          acc.Address,
			Seed:             acc.Seed,
		}
		subAccs = append(subAccs, acc)

	}
	return subAccs, err
}

func parseCoins(coins []string) (float64, error) {
	coin := coins[0]
	amtStr := strings.Replace(coin, constants.Denom, "", -1)
	amtFloat, err := strconv.ParseFloat(amtStr, 64)

	if err != nil {
		return float64(0), nil
	}
	return amtFloat, nil
}

func parseFloat64ToStr(amt float64) string {
	return strconv.FormatFloat(amt, 'f', -1, 64)
}
