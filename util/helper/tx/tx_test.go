package tx

import (
	"testing"

	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"os"
)

func TestMain(m *testing.M) {
	conf.NodeUrl = "http://localhost:1317"
	conf.ChainId = "rainbow-dev"

	conf.FaucetSeed = "cube water sing thunder rib buyer assume rebuild cigar earn slight canoe apart grocery image satisfy genre woman mother can client science this tag"

	conf.BlockInterval = 5
	conf.DefaultReceiverAddr = "faa1r5q5wqwctgfpt3p56qsctptrcq4st6lssyzx65"

	code := m.Run()
	os.Exit(code)
}

func TestSendTransferTx(t *testing.T) {
	type args struct {
		senderInfo   types.AccountInfo
		receiver     string
		generateOnly bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test send transfer tx",
			args: args{
				senderInfo: types.AccountInfo{
					LocalAccountName: "mock-faucet",
					Password:         "1234567890",
					AccountNumber:    "23",
					Sequence:         "1",
					Address:          "faa1mhx2fgwds8uszeazl3au6r0xceppj9xrxavpud",
				},
				receiver:     conf.DefaultReceiverAddr,
				generateOnly: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := SendTransferTx(tt.args.senderInfo, tt.args.receiver, "", tt.args.generateOnly)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(res))
		})
	}
}
