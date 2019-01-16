package sign

import (
	"bytes"
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
	"strings"
	"time"
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

	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	Cdc = cdc
}

// sign tx
func signTx(unsignedTx types.TxDataRes, senderInfo types.AccountInfo) ([]byte, error) {
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
		Tx:            unsignedTx.Value,
		Name:          senderInfo.LocalAccountName,
		Password:      senderInfo.Password,
		ChainID:       conf.ChainId,
		AccountNumber: fmt.Sprintf("%d", accountNumber),
		Sequence:      fmt.Sprintf("%d", sequence),
		AppendSig:     true,
	}

	// send sign tx request
	reqBytes, err := json.Marshal(signTxReq)
	//log.Printf("%s\n", reqBytes)
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

func BroadcastTx(txBody string) ([]byte, error) {

	reqBytes := []byte(txBody)

	reqBuffer := bytes.NewBuffer(reqBytes)
	uri := constants.UriTxBroadcast + "?async=true"
	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	// handle response
	if err != nil {
		return nil, err
	}
	if statusCode != constants.StatusCodeOk {
		return nil, fmt.Errorf("unexcepted status code: %v", statusCode)
	}

	if strings.Contains(string(resBytes), "invalid") {
		log.Printf("%s\n", resBytes)
		return nil, fmt.Errorf("unexcepted information: %v", "invalid")
	}

	a := strings.Contains(string(resBytes), "check_tx") && strings.Contains(string(resBytes), "deliver_tx")
	b := strings.Contains(string(resBytes), "hash") && strings.Contains(string(resBytes), "height")
	if a && b {
		return resBytes, nil
	} else {
		log.Printf("check_tx check broadcast information failed\n")
		return nil, fmt.Errorf("check broadcast information failed")
	}
	return resBytes, nil
}

// broadcast many txs for an account, the num is test-num
func BroadcastTxForAccount(SignedDataArray []conf.SignedData, accountIndex int, testDataChan chan types.TestPressData) (int, error) {
	var err error
	var counter int
	var testData types.TestPressData
	testData.AccountIndex = accountIndex
	log.Printf("test account no.%d start %d txs\n", accountIndex, len(SignedDataArray[accountIndex].SignedDataArray))
	for i := 0; i < len(SignedDataArray[accountIndex].SignedDataArray); i++ {
		_, err = BroadcastTx(SignedDataArray[accountIndex].SignedDataArray[i])
		// handle response
		if err != nil {
			log.Printf(err.Error())
			log.Printf("test account no.%d succeed %d txs, sum is %d\n", accountIndex, counter, len(SignedDataArray[accountIndex].SignedDataArray))
			testData.SuccessIndex = counter
			testDataChan <- testData
			return counter, err
		} else {
			counter = counter + 1
		}
	}
	log.Printf("test account no.%d succeed %d txs, sum is %d\n", accountIndex, counter, len(SignedDataArray[accountIndex].SignedDataArray))
	testData.SuccessIndex = counter
	testDataChan <- testData
	return counter, nil
}

// broadcast many txs for an account, the num is test-num
func BroadcastTxForAccountByTime(SignedDataArray []conf.SignedData, accountIndex int, testDataChan chan types.TestPressData) (int, error) {
	var err error
	var counter int
	var testData types.TestPressData
	testData.AccountIndex = accountIndex
	log.Printf("test account no.%d start %d txs\n", accountIndex, len(SignedDataArray[accountIndex].SignedDataArray))
	timeTemp := time.Now()
	for i := 0; i < len(SignedDataArray[accountIndex].SignedDataArray); i++ {
		_, err = BroadcastTx(SignedDataArray[accountIndex].SignedDataArray[i])
		// handle response
		if err != nil {
			log.Printf(err.Error())
			log.Printf("test account no.%d succeed %d txs, sum is %d\n", accountIndex, counter, len(SignedDataArray[accountIndex].SignedDataArray))
			testData.SuccessIndex = counter
			testDataChan <- testData
			return counter, err
		} else {
			counter = counter + 1
		}

		if i%2500 == 0 {
			if timeTemp.Add(time.Second * 5).After(time.Now()) {
				time.Sleep(timeTemp.Add(time.Second * 5).Sub(time.Now()))
				log.Printf("test broadcast %d\n", i)
			}
			timeTemp = time.Now()
		}
	}
	log.Printf("test account no.%d succeed %d txs, sum is %d\n", accountIndex, counter, len(SignedDataArray[accountIndex].SignedDataArray))
	testData.SuccessIndex = counter
	testDataChan <- testData
	return counter, nil
}

// generate signed tx
func GenSignedTxData(senderInfo types.AccountInfo, receiver string, resChan chan types.GenSignedTxDataRes, chanNum int) {
	var (
		unsignedTx, signedTx auth.StdTx
		uu                   types.TxDataRes
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
	signedTxBytes, err := signTx(uu, senderInfo)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}
	err = Cdc.UnmarshalJSON(signedTxBytes, &signedTx)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}

	/*	// build signed data
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
			Fee: types.StdFee{
				Amount: signedTx.Fee.Amount,
				Gas:    uint64(signedTx.Fee.Gas),
			},
			Signatures: []types.StdSignature{stdSign},
			Memo:       signedTx.Memo,
		}

		postTxBytes, err := json.Marshal(postTx)

		if err != nil {
			log.Printf("%v: cdc marshal json fail: %v\n", method, err)
			return
		}

		signedTxDataRes.Res = string(postTxBytes)*/

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

// generate signed tx
func GenSignedTxDataFromSingleFaucet(faucetAddr string, senderInfo types.AccountInfo, receiver string, resChan chan types.GenSignedTxDataRes, chanNum int) {
	var (
		unsignedTx, signedTx types.TxDataRes
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
	unsignedTxBytes, err := tx.SendTransferTxFromFaucet(senderInfo, faucetAddr, "0.001iris", true)
	if err != nil {
		log.Printf("%v: build unsigned tx failed: %v\n", method, err)
		return
	}
	//log.Printf("%s\n", unsignedTxBytes)
	err = json.Unmarshal(unsignedTxBytes, &unsignedTx)
	//err = Cdc.UnmarshalJSON(unsignedTxBytes, &unsignedTx)
	if err != nil {
		log.Printf("%v: build unsigned tx failed: %v\n", method, err)
		return
	}

	// sign tx
	//fmt.Println(unsignedTx)
	signedTxBytes, err := signTx(unsignedTx, senderInfo)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}
	//log.Printf("%s\n", signedTxBytes)

	err = json.Unmarshal(signedTxBytes, &signedTx)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}

	postTx := types.TxBroadcast{
		Tx: signedTx.Value,
	}

	postTxBytes, err := json.Marshal(postTx)
	//log.Printf("%s\n", postTxBytes)
	if err != nil {
		log.Printf("%v: cdc marshal json fail: %v\n", method, err)
		return
	}

	signedTxDataRes.Res = string(postTxBytes)
}

func GenSignedTxDataByAmountAndFaucet(amount string, faucetAddr string, senderInfo types.AccountInfo, receiver string, resChan chan types.GenSignedTxDataRes, chanNum int) {
	var (
		unsignedTx, signedTx types.TxDataRes
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
	unsignedTxBytes, err := tx.SendTransferTxFromFaucet(senderInfo, faucetAddr, amount, true)
	if err != nil {
		log.Printf("%v: build unsigned tx failed: %v\n", method, err)
		return
	}
	//log.Printf("%s\n", unsignedTxBytes)
	err = json.Unmarshal(unsignedTxBytes, &unsignedTx)
	//err = Cdc.UnmarshalJSON(unsignedTxBytes, &unsignedTx)
	if err != nil {
		log.Printf("%v: build unsigned tx failed: %v\n", method, err)
		return
	}

	// sign tx
	//fmt.Println(unsignedTx)
	signedTxBytes, err := signTx(unsignedTx, senderInfo)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}
	//log.Printf("%s\n", signedTxBytes)

	err = json.Unmarshal(signedTxBytes, &signedTx)
	if err != nil {
		log.Printf("%v: sign tx failed: %v\n", method, err)
		return
	}

	postTx := types.TxBroadcast{
		Tx: signedTx.Value,
	}

	postTxBytes, err := json.Marshal(postTx)
	//log.Printf("%s\n", postTxBytes)
	if err != nil {
		log.Printf("%v: cdc marshal json fail: %v\n", method, err)
		return
	}

	signedTxDataRes.Res = string(postTxBytes)

}
