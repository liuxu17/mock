package sign

import (
	"encoding/json"
	"errors"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/crypto/keys/hd"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	sdk "github.com/irisnet/irishub/types"
	. "github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
	"log"
	"strconv"
)

type StdSignMsg struct {
	ChainID       string      `json:"chain_id"`
	AccountNumber uint64      `json:"account_number"`
	Sequence      uint64      `json:"sequence"`
	Fee           auth.StdFee `json:"fee"`
	Msgs          []sdk.Msg   `json:"msgs"`
	Memo          string      `json:"memo"`
}

// get message bytes
func (msg StdSignMsg) Bytes() []byte {
	return auth.StdSignBytes(msg.ChainID, msg.AccountNumber, msg.Sequence, msg.Fee, msg.Msgs, msg.Memo)
}

const (
	amtV    = "1000000000000000"
	feeAmtV = "10000000000000000"
	denom   = "iris-atto"
	gas     = uint64(20000)
	memo    = ""
	seed    = "tube lonely pause spring gym veteran know want grid tired taxi such same mesh charge orient bracket ozone concert once good quick dry boss"
)

func InitAccountSignProcess(fromAddr string, mnemonic string) (types.AccountTestPrivateInfo, error) {
	var Account types.AccountTestPrivateInfo

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return Account, err
	}
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)
	if err != nil {
		return Account, err
	}
	pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()

	acc, err := account.GetAccountInfo(fromAddr)
	sequence, err := strconv.Atoi(acc.Sequence)
	if err != nil {
		return Account, err
	}
	accountNumber, err := strconv.Atoi(acc.AccountNumber)
	if err != nil {
		return Account, err
	}

	Account.PrivateKey = derivedPriv
	Account.PubKey = pubk.Bytes()
	Account.Addr = fromAddr
	Account.AccountNumber = uint64(accountNumber)
	Account.Sequence = uint64(sequence)
	return Account, err
}

func GenSignTxByTend(testNum int, chainId string, subFaucets []SubFaucet, accountPrivate types.AccountTestPrivateInfo) ([]string, error) {

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		"tendermint/PubKeySecp256k1", nil)

	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterInterface((*sdk.Tx)(nil), nil)

	from, err := sdk.AccAddressFromBech32(subFaucets[testNum].FaucetAddr)

	if err != nil {
		return nil, errors.New("err in address to String")
	}
	amount, ok := sdk.NewIntFromString(amtV)
	coins := sdk.Coins{{Denom: denom, Amount: amount}}

	input := bank.Input{Address: from, Coins: coins}


	feea, ok := sdk.NewIntFromString(feeAmtV)
	feeAmt := sdk.Coins{{Denom: denom, Amount: feea}}
	if !ok {
		return nil, errors.New("err in String to int")
	}
	var msgs []sdk.Msg

	fee := auth.StdFee{Amount: feeAmt, Gas: gas}

	sigMsg := StdSignMsg{
		ChainID:       chainId,
		AccountNumber: accountPrivate.AccountNumber,
		Sequence:      accountPrivate.Sequence,
		Memo:          "",
		Msgs:          msgs,
		Fee:           fee,
	}
	var signedData []string
	priv := secp256k1.PrivKeySecp256k1(accountPrivate.PrivateKey)
	sigChan := make(chan auth.StdTx, 30)
	counter := 0
	for i := 0; i < testNum; i++ {
		to, err := sdk.AccAddressFromBech32(subFaucets[counter].FaucetAddr)
		if err != nil {
			return nil, errors.New("err 2 in address to String")
		}
		output := bank.Output{Address: to, Coins: coins}
		msgs = []sdk.Msg{bank.MsgSend{
			Inputs:  []bank.Input{input},
			Outputs: []bank.Output{output},
		}}

		go genSignedDataByTend(priv, sigMsg, accountPrivate, sigChan, msgs, fee)
		tx := <- sigChan
		bz, _ := cdc.MarshalJSON(tx)
		var signedTx types.TxDataRes
		err = json.Unmarshal(bz, &signedTx)
		if err != nil {
			log.Printf("%v: sign tx failed: %v\n", "SignByTend", err)
			return nil, err
		}

		postTx := types.TxBroadcast{
			Tx: signedTx.Value,
		}
		postTx.Tx.Msgs[0].Type = "irishub/bank/Send"
		postTxBytes, err := json.Marshal(postTx)
		//log.Printf("%s\n", postTxBytes)
		if err != nil {
			log.Printf("%v: cdc marshal json fail: %v\n", "", err)
			return nil, err
		}
		signedData = append(signedData, string(postTxBytes))
		sigMsg.Sequence = sigMsg.Sequence + 1
		counter = counter + 1
		if counter == len(subFaucets) {
			counter = 0
		}
	}
	return signedData, nil
}

func genSignedDataByTend(priv secp256k1.PrivKeySecp256k1, sigMsg StdSignMsg, accountPrivate types.AccountTestPrivateInfo, sigChan chan auth.StdTx, msgs []sdk.Msg, fee auth.StdFee){
	sigBz := sigMsg.Bytes()
	sigByte, err := priv.Sign(sigBz)
	if err != nil {
		return
	}
	sig := auth.StdSignature{
		PubKey:        priv.PubKey(),
		Signature:     sigByte,
		AccountNumber: accountPrivate.AccountNumber,
		Sequence:      sigMsg.Sequence,
	}

	sigs := []auth.StdSignature{sig}
	tx := auth.NewStdTx(msgs, fee, sigs, memo)
	sigChan <- tx
}