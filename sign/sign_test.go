package sign

import (
	"testing"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/constants"
		)

func TestBroadcastSignedTx(t *testing.T) {

	type args struct {
		senderInfo types.AccountInfo
		receiver   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test broadcast a signed tx",
			args: args{
				senderInfo: types.AccountInfo{
					LocalAccountName: constants.MockFaucetName,
					Password:         constants.MockFaucetPassword,
					AccountNumber:    "0",
					Sequence:         "3",
					Address:          conf.FaucetAddress,
				},
				receiver: "faa1z75mnqnzkr72ehmqh2zcx38fmn52af8sk6rwx5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GenSignedTxData(tt.args.senderInfo, tt.args.receiver)
			if err != nil {

			}
			t.Log(string(res))
		})
	}
}
