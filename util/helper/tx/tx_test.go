package tx

import (
	"testing"

	"github.com/kaifei-bianjie/mock/types"
)

func TestSendTransferTx(t *testing.T) {
	type args struct {
		senderInfo   types.AccountInfo
		receiver     string
		generateOnly bool
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "test send transfer tx",
			args: args{
				senderInfo: types.AccountInfo{
					LocalAccountName: "mock-faucet",
					Password: "1234567890",
					AccountNumber: "0",
					Sequence: "1",
					Address: "faa1q5nlka2hwqs86e92704tng5u0tpq700mpwx6l2",
				},
				receiver: "faa1dxzkswsdvc3r0jky0388shcd329hyrpcryq40w",
				generateOnly: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendTransferTx(tt.args.senderInfo, tt.args.receiver, tt.args.generateOnly); err != nil {
				t.Fatal(err)
			}
		})
	}
}
