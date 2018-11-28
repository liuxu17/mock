package key

import (
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"github.com/kaifei-bianjie/mock/util/helper/tx"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/conf"
	"time"
	"log"
	"github.com/kaifei-bianjie/mock/util/helper"
	"fmt"
)

// create account and return account info
func CreateAccounts(num int) ([]types.AccountInfo, error) {
	var (
		successCreatedAccs, ownTokenAccs, accountsInfo []types.AccountInfo
		faucetSequence                                 int64
		method                                         = "CreateAccouts"
	)

	createKeyChan := make(chan types.AccountInfo)
	accInfoChan := make(chan types.AccountInfo)

	// get faucet info
	// sequence of faucet will increment in loop
	faucetInfo := types.AccountInfo{
		LocalAccountName: constants.MockFaucetName,
		Password:         constants.MockFaucetPassword,
		Address:          conf.MockFaucetAddress,
	}
	faucetAccount, err := account.GetAccountInfo(faucetInfo.Address)
	if err != nil {
		log.Printf("%v: get faucet info fail: %v\n", method, err)
		return nil, err
	}
	faucetInfo.AccountNumber = faucetAccount.AccountNumber
	faucetSequence, err = helper.ConvertStrToInt64(faucetAccount.Sequence)
	if err != nil {
		log.Printf("%v: convert sequence to int64 fail: %v\n", method, err)
		return nil, err
	}

	// create account
	createKey := func(index int, accChan chan types.AccountInfo) {
		var accountInfo types.AccountInfo
		keyName := account.GenKeyName(constants.KeyNamePrefix, index)

		// create account
		address, err := account.CreateAccount(keyName, constants.KeyPassword, "")
		if err != nil {
			log.Printf("%v: create key fail: %v\n", method, err)
			accChan <- accountInfo
		}
		log.Printf("%v: account which name is %v create success\n",
			method, keyName)

		accountInfo.LocalAccountName = keyName
		accountInfo.Password = constants.KeyPassword
		accountInfo.Address = address

		accChan <- accountInfo
	}

	// receive token from faucet
	receiveToken := func(address string, accInfo types.AccountInfo) (types.AccountInfo, error) {
		// receive token from faucet
		// get faucet account info
		senderInfo := faucetInfo
		senderInfo.Sequence = fmt.Sprintf("%v", faucetSequence)

		// faucet transfer token
		_, err = tx.SendTransferTx(senderInfo, address, "", false)
		if err != nil {
			log.Printf("%v: faucet transfer token to %v fail: %v\n",
				method, accInfo.LocalAccountName, err)
			return accInfo, err
		}
		faucetSequence += 1
		log.Printf("%v: faucet transfer token to %v success\n",
			method, accInfo.LocalAccountName)
		return accInfo, nil
	}

	// get account info which created: account_number, sequence
	getAccountInfo := func(accInfo types.AccountInfo, accInfoChan chan types.AccountInfo) {
		// get account info
		acc, err := account.GetAccountInfo(accInfo.Address)
		if err != nil {
			log.Printf("%v: get %v info fail: %v\n",
				method, accInfo.LocalAccountName, err)
			accInfoChan <- accInfo
		}
		accInfo.AccountNumber = acc.AccountNumber
		accInfo.Sequence = acc.Sequence
		accInfoChan <- accInfo
	}

	// use goroutine to create account
	for i := 1; i <= num; i++ {
		go createKey(i, createKeyChan)
	}

	counter := 0
	for {
		accInfo := <-createKeyChan
		if accInfo.Address != "" {
			successCreatedAccs = append(successCreatedAccs, accInfo)
		}
		counter ++
		if counter == num {
			log.Printf("%v: all create key goroutine over\n", method)
			log.Printf("%v: except create %v accounts, but success create %v accounts",
				method, num, len(successCreatedAccs))
			break
		}
	}

	// loop transfer token to acc
	// can't use goroutine because of sequence in tx must in order,
	// for example, tx which sequence is 35 shouldn't be broadcasted to blockchain
	// while tx which sequence is 34 hasn't be broadcasted to blockchain
	for _, acc := range successCreatedAccs {
		ownTokenAcc, err := receiveToken(acc.Address, acc)
		if err != nil {

		}
		ownTokenAccs = append(ownTokenAccs, ownTokenAcc)
	}

	// note: can't get account info if not wait 2 block
	log.Printf("%v: sleep %vs before get account sequence\n",
		method,conf.BlockInterval*2)
	time.Sleep(time.Second * time.Duration(conf.BlockInterval*2))
	log.Printf("%v: sleep over\n", method)

	// use goroutine to get sequence of account
	for _, acc := range ownTokenAccs {
		go getAccountInfo(acc, accInfoChan)
	}

	counter = 0
	for {
		accInfo := <-accInfoChan
		if accInfo.AccountNumber != "" {
			accountsInfo = append(accountsInfo, accInfo)
		}
		counter ++

		if counter == len(ownTokenAccs) {
			break
		}
	}

	return accountsInfo, nil
}
