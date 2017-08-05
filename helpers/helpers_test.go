package helpers

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestStringInSlice(t *testing.T) {
	type args struct {
		search string
		sl     []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"string is in slice",
			args{
				"a",
				[]string{"c", "d", "a"}},
			true,
		},
		{
			"string is not in slice",
			args{
				"ab",
				[]string{"c", "d", "a"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.search, tt.args.sl); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "helpertest")
	defer os.Remove(tmpfile.Name())
	type args struct {
		fpath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"file exists",
			args{tmpfile.Name()},
			true,
		},
		{
			"file not exists",
			args{"notexists"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.fpath); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
