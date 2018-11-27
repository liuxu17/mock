package types

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
