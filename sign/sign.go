package sign

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	sdk "github.com/irisnet/irishub/types"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/kaifei-bianjie/mock/util/helper/tx"
	"log"
)

const (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "faa"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = "fap"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = "fva"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = "fvp"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = "fca"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = "fcp"
)

var (
	Cdc *codec.Codec
)

//custom tx codec
func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()

	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	Cdc = cdc
}

// sign tx
func signTx(unsignedTx auth.StdTx, senderInfo types.AccountInfo) ([]byte, error) {
	// build request
	accountNumber, err := helper.ConvertStrToInt64(senderInfo.AccountNumber)
	if err != nil {
		return nil, err
	}
	sequence, err := helper.ConvertStrToInt64(senderInfo.Sequence)
	if err != nil {
		return nil, err
	}
	signTxReq := types.SignTxReq{
		Tx:            unsignedTx,
		Name:          senderInfo.LocalAccountName,
		Password:      senderInfo.Password,
		ChainID:       conf.ChainId,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		AppendSig:     true,
	}

	// send sign tx request
	reqBytes, err := Cdc.MarshalJSON(signTxReq)
	if err != nil {
		return nil, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)
	statusCode, resBytes, err := helper.HttpClientPostJsonData(constants.UriTxSign, reqBuffer)

	// handle response
	if err != nil {
		return nil, err
	}

	if statusCode != constants.StatusCodeOk {
		return nil, fmt.Errorf("unexcepted status code: %v", statusCode)
	}

	return resBytes, nil
}

// generate signed tx
func GenSignedTxData(senderInfo types.AccountInfo, receiver string, resChan chan types.GenSignedTxDataRes, chanNum int) {
	var (
		unsignedTx, signedTx auth.StdTx
		method               = "GenSignedTxData"
	)
	log.Printf("%v: %v goroutine begin gen signed data\n", method, chanNum)

	signedTxDataRes := types.GenSignedTxDataRes{
		ChanNum: chanNum,
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v: failed: %v\n", method, err)
		}

		//log.Printf("%v: signed tx data: %v\n", method, signedTxDataRes.Res)
		resChan <- signedTxDataRes
	}()

	// build unsigned tx
	unsignedTxBytes, err := tx.SendTransferTx(senderInfo, receiver, "0.01iris", true)
	if err != nil {
		log.Printf("%v: build unsigned tx failed: %v\n", method, err)
		return
	}
	err = Cdc.UnmarshalJSON(unsignedTxBytes, &unsignedTx)
	if err != nil {
		log.Printf("%v: build unsigned tx failed: %v\n", method, err)
		return
	}

	// sign tx
	signedTxBytes, err := signTx(unsignedTx, senderInfo)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}
	err = Cdc.UnmarshalJSON(signedTxBytes, &signedTx)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}

	// build signed data
	msgBytes, err := Cdc.MarshalJSON(signedTx.Msgs[0])
	if err != nil {
		log.Printf("%v: build post tx data failed: %v\n", method, err)
		return
	}

	signature := signedTx.Signatures[0]

	stdSign := types.StdSignature{
		PubKey:        signature.PubKey.Bytes(),
		Signature:     signature.Signature,
		AccountNumber: signature.AccountNumber,
		Sequence:      signature.Sequence,
	}

	postTx := types.PostTx{
		Msgs: []string{string(msgBytes)},
		Fee: auth.StdFee{
			Amount: signedTx.Fee.Amount,
			Gas:    int64(signedTx.Fee.Gas),
		},
		Signatures: []types.StdSignature{stdSign},
		Memo:       signedTx.Memo,
	}

	postTxBytes, err := json.Marshal(postTx)

	if err != nil {
		log.Printf("%v: cdc marshal json fail: %v\n", method, err)
		return
	}
	signedTxDataRes.Res = base64.StdEncoding.EncodeToString(postTxBytes)

	//if err != nil {
	//	log.Printf("broadcast tx failed: %v\n", err)
	//	return nil, err
	//}
	//reqBuffer := bytes.NewBuffer(reqBytes)
	//
	//statusCode, resBytes, err := helper.HttpClientPostJsonData(constants.UriTxBroadcastTx, reqBuffer)
	//
	//if err != nil {
	//	log.Printf("broadcast tx failed: %v\n", err)
	//	return nil, err
	//}
	//
	//if statusCode != constants.StatusCodeOk {
	//	log.Printf("broadcast tx failed, unexcepted status code: %v\n", statusCode)
	//	return nil, err
	//}
	//
	//return resBytes, nil
}
