package types

import (
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/types"
)

type AccountInfo struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	AccountNumber    string `json:"account_number"`
	Sequence         string `json:"sequence"`
	Address          string `json:"address"`
}

type AccountInfoRes struct {
	AccountNumber string   `json:"account_number"`
	Sequence      string   `json:"sequence"`
	Address       string   `json:"address"`
	Coins         []string `json:"coins"`
}

type BaseTx struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	AccountNumber    string `json:"account_number"`
	Sequence         string `json:"sequence"`
	Gas              string `json:"gas"`
	Fees             string `json:"fee"`
	Memo             string `json:"memo"`
}

type TransferTxReq struct {
	Amount string `json:"amount"`
	Sender string `json:"sender"`
	BaseTx BaseTx `json:"base_tx"`
}

type ErrorRes struct {
	RestAPI      string `json:"rest api"`
	Code         int    `json:"code"`
	ErrorMessage string `json:"err message"`
}

type KeyCreateReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

type KeyCreateRes struct {
	Address string `json:"address"`
}

type SignTxReq struct {
	Tx            PostTx `json:"tx"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainID       string `json:"chain_id"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
	AppendSig     bool   `json:"append_sig"`
}

type PostTxReq struct {
	Tx auth.StdTx `json:"tx"`
}

type GenSignedTxDataRes struct {
	Res     string
	ChanNum int
}

type PostTx struct {
	Msgs       []TxDataInfo   `json:"msg"`
	Fee        StdFee         `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

type StdFee struct {
	Amount types.Coins `json:"amount"`
	Gas    string      `json:"gas"`
}

type StdSignature struct {
	PubKey        PubKey `json:"pub_key"` // optional
	Signature     string `json:"signature"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

type PubKey struct {
	Type  string `json:"type"` // optional
	Value string `json:"value"`
}

type KeyInfo struct {
	PubKey  string `json:"pub_key"`
	Address string `json:"address"`
	Name    string `json:"name"`
	KeyType string `json:"type"`
}

type TxDataRes struct {
	Type  string `json:"type"`
	Value PostTx `json:"value"`
}

type TxDataInfo struct {
	Type  string      `json:"type"`
	Value TxDataValue `json:"value"`
}

type TxDataValue struct {
	Input  []InOutPutData `json:"inputs"`
	Output []InOutPutData `json:"outputs"`
}

type InOutPutData struct {
	Address string      `json:"address"`
	Amount  types.Coins `json:"coins"`
}

type TxBroadcast struct {
	Tx PostTx `json:"tx"`
}
