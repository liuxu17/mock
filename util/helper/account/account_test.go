package account

import (
	"testing"

			"encoding/json"
	"github.com/kaifei-bianjie/mock/util/contants"
	"math/rand"
)

func TestGenKeyName(t *testing.T) {
	type args struct {
		namePrefix string
		num        int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenKeyName(tt.args.namePrefix, tt.args.num); got != tt.want {
				t.Errorf("GenKeyName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	type args struct {
		name     string
		password string
		seed     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test create account",
			args: args{
				name:     GenKeyName(contants.KeyNamePrefix, rand.Intn(10)),
				password: "1234567890",
				seed:     "",
			},
		},
		//{
		//	name: "test recover account",
		//	args: args{
		//		name:     "mock-faucet",
		//		password: "1234567890",
		//		seed: "nephew pupil few cash liberty sorry stay brand east antenna civil cat area endorse wheel chronic inform diesel next drip style neither salad nominee",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CreateAccount(tt.args.name, tt.args.password, tt.args.seed)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("address of new account is %v\n", res)
		})
	}
}

func TestGetAccountInfo(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test get account info",
			args: args{
				address: "faa1q5nlka2hwqs86e92704tng5u0tpq700mpwx6l2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetAccountInfo(tt.args.address)
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
