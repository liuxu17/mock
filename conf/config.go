package conf

import "os"

var (
	NodeUrl string
	ChainId string

	FaucetSeed    string
	FaucetAddress string

	BlockInterval       int
	DefaultReceiverAddr string

	DefaultHome = os.ExpandEnv("$HOME") + "/mock"

	KeyFaucetSeed    = "faucet_seed"
	KeyFaucetAddress = "faucet_address"
)

//func init()  {
//	NodeUrl = "http://localhost:1317"
//	ChainId = "rainbow-dev"
//	BlockInterval = 5
//	FaucetSeed = "recycle light kid spider fire disorder relax end stool hip child leaf wild next veteran start theory pretty salt rich avocado card enact april"
//	FaucetAddress = "faa1t5wlur60xzzcxpgjn0d5y8ge7fsdmp7jejl7am"
//	DefaultReceiverAddr = FaucetAddress
//}
