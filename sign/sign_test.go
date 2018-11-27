package sign

import (
	"testing"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenSignData(t *testing.T) {
	type args struct {
		m sendBody
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "test get sign data",
			args: args{
				m: sendBody{
					LocalAccountName: "mock-faucet",
					Password: "1234567890",
					ChainID: conf.MockChainId,
					AccountNumber: 0,
					Sequence: 28,
					Gas: 200000,
					Receiver: "faa1waty0rrpsww3wrwyxhxu0memtr74ar5pegk0zc",
					Amount: []types.Coin{
						{
							Denom: "iris-atto",
							Amount: sdk.NewIntWithDecimal(2, 18),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GenSignData(tt.args.m)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		})
	}
}
