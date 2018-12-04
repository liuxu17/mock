package service

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var (
	subFaucets []conf.SubFaucet
)

func TestMain(m *testing.M) {
	// set config var
	conf.NodeUrl = "http://localhost:1317"
	conf.ChainId = "rainbow-dev"
	conf.BlockInterval = 5
	conf.DefaultReceiverAddr = "faa1jyj90se9mel2smn3vr4u9gzg03acwuy8h44q3m"

	// read config from config file
	// check config file if exists
	confHome := conf.DefaultHome
	confFilePath := fmt.Sprintf("%v/%v", confHome, constants.ConfigFileName)

	exists, err := helper.CheckFileExist(confFilePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(fmt.Errorf("can't find config file in %v\n", confFilePath))
	}

	// read config from config file
	viper.SetConfigName("config")
	viper.AddConfigPath(confHome)
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.UnmarshalKey("sub_faucets", &subFaucets)
	if err != nil {
		panic(err)
	}
	if len(subFaucets) <= 0 {
		panic(fmt.Errorf("can't read sub_faucets config"))
	}

	existCode := m.Run()

	os.Exit(existCode)
}

func TestBatchGenSignedTxData(t *testing.T) {
	type args struct {
		num        int
		subFaucets []conf.SubFaucet
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test batch gen signed tx",
			args: args{
				num:        30,
				subFaucets: subFaucets,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := BatchGenSignedTxData(tt.args.num, tt.args.subFaucets)
			t.Log(res)
		})
	}
}
