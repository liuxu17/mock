package key

import (
	"github.com/kaifei-bianjie/mock/util/contants"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"github.com/kaifei-bianjie/mock/util/helper/tx"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/conf"
	"time"
)

// create account and return account info
func CreateAccounts(num int) ([]types.AccountInfo, error) {
	var accountsInfo []types.AccountInfo
	faucetInfo := types.AccountInfo{
		LocalAccountName: conf.MockFaucetName,
		Password: conf.MockFaucetPassword,
		Address: conf.MockFaucetAddress,
	}

	// TODO: use goroutine do these task
	for i := 1; i <= num ; i++  {
		keyName := account.GenKeyName(constants.KeyNamePrefix, i)
		accountInfo := types.AccountInfo{
			LocalAccountName: keyName,
			Password:         constants.KeyPassword,
		}

		// create account
		address, err := account.CreateAccount(keyName, constants.KeyPassword, "")
		if err != nil {
			return accountsInfo, err
		}

		// get token from faucet
		// get senderInfo info(faucet account info)
		senderInfo := faucetInfo
		faucetAccount, err := account.GetAccountInfo(senderInfo.Address)
		if err != nil {
			return accountsInfo, err
		}
		senderInfo.AccountNumber = faucetAccount.AccountNumber
		senderInfo.Sequence = faucetAccount.Sequence

		// faucet transfer token
		err = tx.SendTransferTx(senderInfo, address, false)
		if err != nil {
			return accountsInfo, err
		}

		// note: can't get account info if not wait 2 block
		time.Sleep(time.Second * time.Duration(conf.BlockInterval * 2))

		// get account info
		acc, err := account.GetAccountInfo(address)
		accountInfo.Address = address
		accountInfo.AccountNumber = acc.AccountNumber
		accountInfo.Sequence = acc.Sequence

		accountsInfo = append(accountsInfo, accountInfo)
	}

	return accountsInfo, nil
}
