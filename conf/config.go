package conf

var (
	BlockChainNodeServerUrl string
	MockChainId string
	BlockInterval int
	MockFaucetSeed string
	MockFaucetAddress string
	MockReceiverAddr string
)


func init()  {
	BlockChainNodeServerUrl = "http://localhost:1317"
	MockChainId = "rainbow-dev"
	BlockInterval = 5
	MockFaucetSeed = "recycle light kid spider fire disorder relax end stool hip child leaf wild next veteran start theory pretty salt rich avocado card enact april"
	MockFaucetAddress = "faa1jyj90se9mel2smn3vr4u9gzg03acwuy8h44q3m"
	MockReceiverAddr = MockFaucetAddress
}
