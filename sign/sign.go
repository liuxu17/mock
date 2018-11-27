package sign

//var (
//	Cdc *wire.Codec
//)
//
// custom tx codec
//func init() {
//
//	var cdc = wire.NewCodec()
//	auth.RegisterWire(cdc)
//	bank.RegisterWire(cdc)
//	sdk.RegisterWire(cdc)
//	wire.RegisterCrypto(cdc)
//	Cdc = cdc
//}

type sendBody struct {
	LocalAccountName string    `json:"name"`
	Password         string    `json:"password"`
	ChainID          string    `json:"chain_id"`
	AccountNumber    int64     `json:"account_number"`
	Sequence         int64     `json:"sequence"`
	Gas              int64     `json:"gas"`
	GasAdjustment string `json:"gas_adjustment"`
	Fee           string `json:"fee"`
	Receiver         string
}

//func GenSignData(m sendBody) (string, error) {
//
//}
