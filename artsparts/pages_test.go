package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var mdtext = `newpage: firstpage
# First Page

This ist the text of the first page.
Second line.

newpage: secondpage
# Second Page
This is the second page.`

var mdtextParsed = []Page{
	Page{
		"firstpage",
		"First Page",
		"\nThis ist the text of the first page.\nSecond line.\n\n",
	},
	Page{
		"secondpage",
		"Second Page",
		"This is the second page.\n",
	},
}

func Test_parsePages(t *testing.T) {

	buf := bytes.NewBufferString(mdtext)
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want []Page
	}{
		{
			"parse pages",
			args{buf},
			mdtextParsed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePages(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePages() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_getPage(t *testing.T) {
	type args struct {
		pagePath string
		pages    []Page
	}
	tests := []struct {
		name string
		args args
		want Page
	}{
		{
			"find a page",
			args{
				"secondpage",
				mdtextParsed,
			},
			mdtextParsed[1],
		},
		{
			"find no page",
			args{
				"notexist",
				mdtextParsed,
			},
			Page{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPage(tt.args.pagePath, tt.args.pages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
