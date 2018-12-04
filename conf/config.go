package conf

import "os"

var (
	NodeUrl string
	ChainId string

	FaucetSeed    string
	FaucetName    string
	FaucetAddress string

	BlockInterval       int
	DefaultReceiverAddr string

	DefaultHome = os.ExpandEnv("$HOME") + "/.mock/config"
)

type ConfigContent struct {
	FaucetSeed string      `json:"faucet_seed" mapstructure:"faucet_seed"`
	FaucetName string      `json:"faucet_name" mapstructure:"faucet_name"`
	FaucetAddr string      `json:"faucet_addr" mapstructure:"faucet_addr"`
	SubFaucets []SubFaucet `json:"sub_faucets" mapstructure:"sub_faucets"`
}

type SubFaucet struct {
	FaucetName     string `json:"faucet_name" mapstructure:"faucet_name"`
	FaucetPassword string `json:"faucet_password" mapstructure:"faucet_password"`
	FaucetAddr     string `json:"faucet_addr" mapstructure:"faucet_addr"`
}

//func init() {
//	NodeUrl = "http://localhost:1317"
//	ChainId = "rainbow-dev"
//	BlockInterval = 5
//	FaucetSeed = "recycle light kid spider fire disorder relax end stool hip child leaf wild next veteran start theory pretty salt rich avocado card enact april"
//	FaucetAddress = "faa1jyj90se9mel2smn3vr4u9gzg03acwuy8h44q3m"
//	DefaultReceiverAddr = FaucetAddress
//}
