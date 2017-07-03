package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateData_AddJS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		td   *TemplateData
		args args
	}{
		{"AddJS", &TemplateData{}, args{"test.js"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.td.AddJS(tt.args.s)
			assert.Contains(t, tt.td.JSFiles, tt.args.s)
		})
	}
}

func TestTemplateData_AddCSS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		td   *TemplateData
		args args
	}{
		{"AddCSS", &TemplateData{}, args{"test.js"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.td.AddCSS(tt.args.s)
			assert.Contains(t, tt.td.CSSFiles, tt.args.s)
		})
	}
}
