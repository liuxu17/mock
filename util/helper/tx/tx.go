package tx

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/util/contants"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/conf"
	"encoding/json"
	"bytes"
	"github.com/kaifei-bianjie/mock/util/helper"
)

// send tokens from mockFaucet account to given address
func SendTransferTx(senderInfo types.AccountInfo, receiver string, generateOnly bool) error {
	uri := fmt.Sprintf(contants.UriTransfer, receiver)
	if generateOnly {
		uri = uri + "?generate-only=true"
	}

	req := types.TransferTxReq{
		Amount: contants.MockTransferAmount,
		Sender: senderInfo.Address,
		BaseTx: types.BaseTx{
			LocalAccountName: senderInfo.LocalAccountName,
			Password:         senderInfo.Password,
			ChainID:          conf.MockChainId,
			AccountNumber:    senderInfo.AccountNumber,
			Sequence:         senderInfo.Sequence,
			Gas:              contants.MockDefaultGas,
			Fees:             contants.MockDefaultFee,
			Memo:             fmt.Sprintf("mock test: transfer token"),
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)

	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBuffer)

	if err != nil {
		return err
	}

	if statusCode == contants.StatusCodeOk {
		return nil
	} else {
		errRes := types.ErrorRes{}
		if err := json.Unmarshal(resBytes, &errRes); err != nil {
			return err
		}
		return fmt.Errorf("err code: %v, err msg: %v", errRes.Code, errRes.ErrorMessage)
	}
}
