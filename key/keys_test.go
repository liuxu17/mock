package key

import (
	"testing"
	"encoding/json"
)

func TestCreateAccounts(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test create accounts",
			args: args{
				num: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CreateAccounts(tt.args.num)
			if err != nil {
				t.Fatal(err)
			}
			resBytes, err := json.MarshalIndent(res, "", "")
			if  err != nil {
				t.Fatal(err)
			}
			t.Log(string(resBytes))
		})
	}
}
