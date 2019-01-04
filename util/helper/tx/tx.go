package tx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
)

// send tokens from mockFaucet account to given address
func SendTransferTx(senderInfo types.AccountInfo, receiver string, amount string, generateOnly bool) ([]byte, error) {
	uri := fmt.Sprintf(constants.UriTransfer, receiver)
	if generateOnly {
		uri = uri + "?generate-only=true"
	}

	if amount == "" {
		amount = constants.MockTransferAmount
	}

	req := types.TransferTxReq{
		Amount: amount,
		Sender: senderInfo.Address,
		BaseTx: types.BaseTx{
			LocalAccountName: senderInfo.LocalAccountName,
			Password:         senderInfo.Password,
			ChainID:          conf.ChainId,
			AccountNumber:    senderInfo.AccountNumber,
			Sequence:         senderInfo.Sequence,
			Gas:              constants.MockDefaultGas,
			Fees:             constants.MockDefaultFee,
			Memo:             fmt.Sprintf("mock test: transfer token"),
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)

	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	if err != nil {
		return nil, err
	}

	if statusCode == constants.StatusCodeOk {
		return resBytes, nil
	} else {
		return nil, fmt.Errorf("transfer token fail: %v\n", string(resBytes))
	}
}

// send tokens from mockFaucet account to given address
func SendTransferTxFromFaucet(senderInfo types.AccountInfo, faucetAddr string, amount string, generateOnly bool) ([]byte, error) {
	uri := fmt.Sprintf(constants.UriTransfer, senderInfo.Address)
	if generateOnly {
		uri = uri + "?generate-only=true"
	}

	if amount == "" {
		amount = constants.MockTransferAmount
	}

	req := types.TransferTxReq{
		Amount: amount,
		Sender: faucetAddr,
		BaseTx: types.BaseTx{
			LocalAccountName: senderInfo.LocalAccountName,
			Password:         senderInfo.Password,
			ChainID:          conf.ChainId,
			AccountNumber:    senderInfo.AccountNumber,
			Sequence:         senderInfo.Sequence,
			Gas:              constants.MockDefaultGas,
			Fees:             constants.MockDefaultFee,
			Memo:             fmt.Sprintf("mock test: transfer token"),
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)

	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	if err != nil {
		return nil, err
	}

	if statusCode == constants.StatusCodeOk {
		return resBytes, nil
	} else {
		return nil, fmt.Errorf("transfer token fail: %v\n", string(resBytes))
	}
}
