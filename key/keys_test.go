package key

import (
	"fmt"
	"testing"
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
				num: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//res, err := NewAccount(tt.args.num)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//resBytes, err := json.MarshalIndent(res, "", "")
			//if err != nil {
			//	t.Fatal(err)
			//}
			//t.Log(string(resBytes))
		})
	}
}

func TestMock(t *testing.T) {
	accs := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "2"}
	threadNum := 5
	interval := len(accs) / threadNum

	for i := 0; i < threadNum; i++ {
		var s, e int
		s = i * interval
		e = s + interval

		if i == threadNum-1 {
			e = len(accs)
		}
		fmt.Printf("threadNum %v, from %v to %v, %v\n", i, s, e, accs[s:e])

	}
}
