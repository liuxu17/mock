package sign

import (
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"log"
	"testing"
)

func TestBroadcastSignedTx(t *testing.T) {

	type args struct {
		senderInfo types.AccountInfo
		receiver   string
		resChan    chan types.GenSignedTxDataRes
		chanNUm    int
	}

	resChannel := make(chan types.GenSignedTxDataRes)

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
					Sequence:         "169",
					Address:          conf.FaucetAddress,
				},
				receiver: "faa1z75mnqnzkr72ehmqh2zcx38fmn52af8sk6rwx5",
				chanNUm:  1,
				resChan:  resChannel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenSignedTxData(tt.args.senderInfo, tt.args.receiver, tt.args.resChan, tt.args.chanNUm)

			res := <-tt.args.resChan
			if res.ChanNum != 0 {
				t.Logf("%v build signed tx data over\n", res.ChanNum)
				t.Log(res.Res)
			}
			log.Println(res.ChanNum)
		})
	}
}
