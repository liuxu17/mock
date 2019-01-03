package key

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"github.com/kaifei-bianjie/mock/util/helper/tx"
	"log"
	"strconv"
	"time"
)

// new account do follow things:
// 1. create key
// 2. faucet transfer token to these accounts
// 3. get account info(account_number and sequence)
func NewAccountSingle(num int, home string, faucet string, faucetAddr string) ([]types.AccountInfo, error) {
	var (
		createdAccs []types.AccountInfo
		//distributedTokenAccs, accountsInfo []types.AccountInfo
		method = "NewAccount"
	)
/*	faucetAddr, err := account.GetAccAddr(faucet)
	if err != nil {
		log.Printf("%v: cannot get %v address fail: %v !!!!!!!!!!!!\n",
			method, faucet, err)
		return nil, err
	}*/

	acc, err := account.GetAccountInfo(faucetAddr)
	if err != nil {
		log.Printf("%v: get %v info fail: %v !!!!!!!!!!!!\n",
			method, faucet, err)
		return nil, err
	}

	createKeyChan := make(chan types.AccountInfo, 100000)

	// create account by cmd
	sequence, err :=strconv.Atoi(acc.Sequence)
	if err != nil {
		return nil, err
	}

	for i := 1; i <= num; i++ {
		keyName := account.GenKeyName(constants.KeyNamePrefix, i)
		CreateKeyByCmd(faucet, strconv.Itoa(sequence), acc.AccountNumber, keyName, createKeyChan, home)
		sequence = sequence + 1
		//go CreateKey(keyName, createKeyChan)
	}

	counter := 0
	for {
		accInfo := <-createKeyChan
		if accInfo.Address != "" {
			createdAccs = append(createdAccs, accInfo)
		}
		counter++
		if counter == num {
			log.Printf("%v: all create key goroutine over\n", method)
			log.Printf("%v: except create %v accounts, successful create %v accounts",
				method, num, len(createdAccs))
			break
		}
	}

	return createdAccs, nil
}

// new account do follow things:
// 1. create key
// 2. faucet transfer token to these accounts
// 3. get account info(account_number and sequence)
func NewAccount(num int, subFaucets []conf.SubFaucet) ([]types.AccountInfo, error) {
	var (
		createdAccs, distributedTokenAccs, accountsInfo []types.AccountInfo
		method                                          = "NewAccount"
	)

	createKeyChan := make(chan types.AccountInfo, 10000)
	distributeChan := make(chan []types.AccountInfo)
	accInfoChan := make(chan types.AccountInfo, 10000)

	// use goroutine to create account
	for i := 1; i <= num; i++ {
		keyName := account.GenKeyName(constants.KeyNamePrefix, i)
		go CreateKey(keyName, createKeyChan)
		//go CreateKey(keyName, createKeyChan)
		if (i % 20) == 0 {
			time.Sleep(time.Second * constants.CreateNewAccountDelaySec)
		}
	}

	counter := 0
	for {
		accInfo := <-createKeyChan
		if accInfo.Address != "" {
			createdAccs = append(createdAccs, accInfo)
		}
		counter++
		if counter == num {
			log.Printf("%v: all create key goroutine over\n", method)
			log.Printf("%v: except create %v accounts, successful create %v accounts",
				method, num, len(createdAccs))
			break
		}
	}

	// loop transfer token to acc
	// can't use goroutine because of sequence in tx must in order,
	// for example, tx which sequence is 35 shouldn't be broadcasted to blockchain
	// while tx which sequence is 34 hasn't be broadcasted to blockchain
	distributeToken := func(senderInfo types.AccountInfo, receiverInfos []types.AccountInfo,
		distributeChan chan []types.AccountInfo) {
		var (
			senderSequence  int64
			distributedAccs []types.AccountInfo
			method          = "distributeToken"
		)
		accInfo, err := account.GetAccountInfo(senderInfo.Address)
		if err != nil {
			log.Printf("%v: err is %v\n", method, err)
			distributeChan <- distributedAccs
		}
		senderInfo.AccountNumber = accInfo.AccountNumber
		senderSequence, err = helper.ConvertStrToInt64(accInfo.Sequence)
		if err != nil {
			log.Printf("%v: err is %v\n", method, err)
			distributeChan <- distributedAccs
		}

		for _, receiver := range receiverInfos {
			senderInfo.Sequence = fmt.Sprintf("%v", senderSequence)
			distributedAcc, err := DistributeToken(senderInfo, receiver, "")
			if err != nil {
				log.Printf("%v: err is %v\n", method, err)
			} else {
				senderSequence += 1
				distributedAccs = append(distributedAccs, distributedAcc)
			}
		}

		distributeChan <- distributedAccs
	}

	// TODO: change master-worker mode
	// use sub faucet account to transfer token
	if len(createdAccs) >= len(subFaucets) {
		eachThreadTask := len(createdAccs) / len(subFaucets)

		log.Printf("%v: %v distribute token task assigned to %v sub faucet\n",
			method, len(createdAccs), len(subFaucets))
		for index, subFaucet := range subFaucets[:] {
			var start, end int
			senderInfo := types.AccountInfo{
				LocalAccountName: subFaucet.FaucetName,
				Password:         subFaucet.FaucetPassword,
				Address:          subFaucet.FaucetAddr,
			}

			start = index * eachThreadTask
			end = start + eachThreadTask

			if index == len(subFaucets)-1 {
				end = len(createdAccs)
			}

			log.Printf("%v: sub faucet %v handler accounts from %v to %v\n",
				method, senderInfo.LocalAccountName, start, end)
			go distributeToken(senderInfo, createdAccs[start:end], distributeChan)
		}

		// get result
		counter := 0
		for {
			res := <-distributeChan
			distributedTokenAccs = append(distributedTokenAccs, res...)
			counter++
			if counter == len(subFaucets) {
				log.Printf("%v: all sub faucet distribute token over\n", method)
				break
			}

		}
	}

	// note: can't get account info if not wait 2 block
	log.Printf("%v: sleep %vs before get account sequence\n",
		method, conf.BlockInterval*2)
	time.Sleep(time.Second * time.Duration(conf.BlockInterval*2))
	log.Printf("%v: sleep over\n", method)

	// use goroutine to get accountInfo
	if len(distributedTokenAccs) >= 1 {
		counter = 0
		for _, acc := range distributedTokenAccs {
			go GetAccountInfo(acc, accInfoChan)
			counter++
			if (counter % 20) == 0 {
				time.Sleep(time.Second * constants.CheckAccountInfoDelaySec)
			}
		}

		counter = 0
		for {
			accInfo := <-accInfoChan
			if accInfo.AccountNumber != "" {
				accountsInfo = append(accountsInfo, accInfo)
			}
			counter++

			if counter == len(distributedTokenAccs) {
				log.Printf("%v: get account info over\n", method)
				break
			}
		}
	}

	return accountsInfo, nil
}

// create key and return accountInfo by channel
func CreateKey(keyName string, accChan chan types.AccountInfo) {
	var (
		accountInfo types.AccountInfo
		method      = "CreateKey"
	)

	// create account
	address, err := account.CreateAccount(keyName, constants.KeyPassword, "")
	if err != nil {
		log.Printf("%v: create key fail: %v\n", method, err)
		return
		//accChan <- accountInfo
	}
	log.Printf("%v: account which name is %v create success\n",
		method, keyName)

	accountInfo.LocalAccountName = keyName
	accountInfo.Password = constants.KeyPassword
	accountInfo.Address = address

	accChan <- accountInfo
}

// create key and return accountInfo by channel
func CreateKeyByCmd(faucetName string, sequence string, AccountNum string, keyName string, accChan chan types.AccountInfo, home string) {
	var (
		accountInfo types.AccountInfo
		method      = "CreateKey"
	)

	// create account
	address, err := account.CreateAccountByCmd(keyName, constants.KeyPassword, home)
	if err != nil {
		log.Printf("%v: create key fail: %v\n", method, err)
		return
		//accChan <- accountInfo
	}
	log.Printf("%v: account which name is %v create success\n",
		method, keyName)

	accountInfo.LocalAccountName = faucetName
	accountInfo.Password = constants.KeyPassword
	accountInfo.Sequence = sequence
	accountInfo.AccountNumber = AccountNum
	accountInfo.Address = address

	accChan <- accountInfo
}

// faucet distribute token to account
func DistributeToken(senderInfo, receiverInfo types.AccountInfo, amount string) (types.AccountInfo, error) {
	var (
		method = "DistributeToken"
	)

	// faucet transfer token
	_, err := tx.SendTransferTx(senderInfo, receiverInfo.Address, amount, false)
	if err != nil {
		log.Printf("%v: %v transfer token to %v fail: %v\n",
			method, senderInfo.LocalAccountName, receiverInfo.LocalAccountName, err)
		return receiverInfo, err
	}
	log.Printf("%v: %v transfer token to %v success\n",
		method, senderInfo.LocalAccountName, receiverInfo.LocalAccountName)
	return receiverInfo, nil
}

// get account info, return account info by channel
func GetAccountInfo(accInfo types.AccountInfo, accInfoChan chan types.AccountInfo) {
	var (
		method = "GetAccountInfo"
	)
	// get account info
	acc, err := account.GetAccountInfo(accInfo.Address)
	if err != nil {
		log.Printf("%v: get %v info fail: %v !!!!!!!!!!!!\n",
			method, accInfo.LocalAccountName, err)
		return
	}
	log.Printf("%v: get %v info success\n", method, accInfo.LocalAccountName)
	accInfo.AccountNumber = acc.AccountNumber
	accInfo.Sequence = acc.Sequence
	accInfoChan <- accInfo
}
