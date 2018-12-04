package key

import (
	"testing"

	"encoding/json"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"os"
	"strings"
)

var (
	faucetName = "mock-faucet-007"
)

func TestMain(m *testing.M) {
	// set config var
	conf.NodeUrl = "http://localhost:1317"
	conf.ChainId = "rainbow-dev"
	conf.BlockInterval = 5
	conf.FaucetSeed = "recycle light kid spider fire disorder relax end stool hip child leaf wild next veteran start theory pretty salt rich avocado card enact april"

	// create faucet account
	addr, err := account.CreateAccount(faucetName, constants.MockFaucetPassword, conf.FaucetSeed)
	if err != nil && !strings.Contains(err.Error(), "acount with name") {
		panic(err)
	}
	conf.FaucetAddress = addr

	existCode := m.Run()

	os.Exit(existCode)
}

func TestCreateFaucetSubAccount(t *testing.T) {
	type args struct {
		faucetName   string
		faucetPasswd string
		faucetAddr   string
		subAccNum    int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test create faucet sub account",
			args: args{
				faucetName:   faucetName,
				faucetPasswd: constants.KeyPassword,
				faucetAddr:   conf.FaucetAddress,
				subAccNum:    10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CreateFaucetSubAccount(tt.args.faucetName, tt.args.faucetPasswd, tt.args.faucetAddr, tt.args.subAccNum)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, err := json.MarshalIndent(res, "", "")
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(resBytes))
		})
	}
}
