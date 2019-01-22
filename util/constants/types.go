package constants

const (
	// config file name
	ConfigFileName = "config.json"

	HeaderContentTypeJson = "application/json"

	// key password, prefix of key name
	KeyNamePrefix = "mock"
	KeyPassword   = "1234567890"

	// http uri
	UriKeyCreate     = "/keys"
	UriAccountInfo   = "/auth/accounts/%v"           // format is /auth/accounts/{address}
	UriKeyInfo       = "/keys/%v"                    // format is /auth/accounts/{address}
	UriTransfer      = "/bank/accounts/%s/transfers" // format is /bank/accounts/{address}/transfers
	UriTxSign        = "/tx/sign"
	UriTxBroadcastTx = "/txs/send"
	UriTxBroadcast   = "/tx/broadcast"

	// http status code
	StatusCodeOk       = 200
	StatusCodeConflict = 409

	//go routine delay time
	CreateNewAccountDelaySec = 4
	CheckAccountInfoDelaySec = 2
	SignTxDelaySec           = 1

	KeysAddCmd  = "iriscli keys add "
	KeysShowCmd = "iriscli keys show "

	//
	MockFaucetName     = "mock-faucet"
	MockFaucetPassword = "1234567890"
	MockTransferAmount = "0.03iris"
	MockDefaultGas     = "20000"
	MockDefaultFee     = "5iris"
	Denom              = "iris"
	FeeAtto            = ""
	transferAmountAtto = ""
)
