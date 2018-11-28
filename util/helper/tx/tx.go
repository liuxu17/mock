package tx

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/conf"
	"encoding/json"
	"bytes"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/kaifei-bianjie/mock/util/constants"
)

// send tokens from mockFaucet account to given address
func SendTransferTx(senderInfo types.AccountInfo, receiver string, generateOnly bool) ([]byte, error) {
	uri := fmt.Sprintf(constants.UriTransfer, receiver)
	if generateOnly {
		uri = uri + "?generate-only=true"
	}

	req := types.TransferTxReq{
		Amount: constants.MockTransferAmount,
		Sender: senderInfo.Address,
		BaseTx: types.BaseTx{
			LocalAccountName: senderInfo.LocalAccountName,
			Password:         senderInfo.Password,
			ChainID:          conf.MockChainId,
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
		errRes := types.ErrorRes{}
		if err := json.Unmarshal(resBytes, &errRes); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("err code: %v, err msg: %v", errRes.Code, errRes.ErrorMessage)
	}
}
