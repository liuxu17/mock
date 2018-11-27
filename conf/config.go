package conf

var (
	BlockChainNodeServerUrl string
	MockChainId string
	BlockInterval int
	MockFaucetSeed string
	MockFaucetAddress string
	MockFaucetName string
	MockFaucetPassword string
)


func init()  {
	BlockChainNodeServerUrl = "http://localhost:1317"
	MockChainId = "rainbow-dev"
	BlockInterval = 5
	MockFaucetSeed = "nephew pupil few cash liberty sorry stay brand east antenna civil cat area endorse wheel chronic inform diesel next drip style neither salad nominee"
	MockFaucetAddress = "faa1q5nlka2hwqs86e92704tng5u0tpq700mpwx6l2"
	MockFaucetName = "mock-faucet"
	MockFaucetPassword = "1234567890"
}
