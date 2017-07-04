package shortlink

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShort(t *testing.T) {
	tmpDBFile, _ := ioutil.TempFile("", "db")
	defer os.Remove(tmpDBFile.Name())
	DbFile = tmpDBFile.Name()
	tests := []struct {
		url      string
		shortURL string
		wantErr  bool
	}{
		{"a/long/id", "/s/1", false},
		{"a/c/longer/url", "/s/2", false},
		{"a/long/id", "/s/1", false},
	}
	for _, tt := range tests {
		gotError := false
		got, err := GetShort(tt.url)
		if err != nil {
			gotError = true
			fmt.Println(err)
		}
		assert.Equal(t, tt.shortURL, got)
		assert.Equal(t, tt.wantErr, gotError)
	}
}
func TestAddLink(t *testing.T) {
	tmpDBFile, _ := ioutil.TempFile("", "db")
	defer os.Remove(tmpDBFile.Name())
	DbFile = tmpDBFile.Name()
	tests := []struct {
		url     string
		expID   int
		wantErr bool
	}{
		{"a/b", 1, false},
		{"a/c", 2, false},
		{"a/b", 1, false},
	}
	for _, tt := range tests {
		gotError := false
		got, err := GetID(tt.url)
		if err != nil {
			gotError = true
			fmt.Println(err)
		}
		assert.Equal(t, tt.expID, got)
		assert.Equal(t, tt.wantErr, gotError)
	}
}
