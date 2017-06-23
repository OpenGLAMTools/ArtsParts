package helpers

import "testing"

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
