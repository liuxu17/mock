package helper

import (
	"os"
	"testing"
)

func TestWriteFile(t *testing.T) {
	type args struct {
		filePath string
		content  []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test write file",
			args: args{
				filePath: "./a.json",
				content:  []byte("This is test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteFile(tt.args.filePath, tt.args.content)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateFolder(t *testing.T) {
	type args struct {
		folderPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test create folder if not exists",
			args: args{
				folderPath: os.ExpandEnv("$HOME") + "/output",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateFolder(tt.args.folderPath)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("create folder successful")
		})
	}
}

func TestConvertStrToInt64(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test convert str to int64",
			args: args{
				s: "149999850.52498939631857518",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ConvertStrToInt64(tt.args.s)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		})
	}
}
