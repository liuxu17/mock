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
					Address: "faa1jyj90se9mel2smn3vr4u9gzg03acwuy8h44q3m",
				},
				receiver: "faa1z75mnqnzkr72ehmqh2zcx38fmn52af8sk6rwx5",
				generateOnly: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := SendTransferTx(tt.args.senderInfo, tt.args.receiver, "", tt.args.generateOnly);
			if  err != nil {
				t.Fatal(err)
			}
			t.Log(string(res))
		})
	}
}
