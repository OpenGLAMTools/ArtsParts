package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

func Test_artsPartsApp_defaultTemplateData(t *testing.T) {
	app := &artsPartsApp{
		getSessionValues: func(r *http.Request) (map[string]string, error) {
			return map[string]string{
				"twitter": "user1",
			}, nil
		},
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		app  *artsPartsApp
		args args
		want templateData
	}{
		{
			"",
			app,
			args{nil},
			templateData{
				JSFiles:  []string{"app.js"},
				CSSFiles: []string{"custom.css"},
				JQuery:   true,
				VueJS:    false,
				User:     "user1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.app.defaultTemplateData(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("artsPartsApp.defaultTemplateData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_artsPartsApp_artwork(t *testing.T) {
	// TODO: test is not complete
	artspartsApp, _ := artsparts.NewApp("../test")
	app := &artsPartsApp{
		artsparts: artspartsApp,
		getSessionValues: func(r *http.Request) (map[string]string, error) {
			return map[string]string{
				"twitter": "user1",
			}, nil
		},
		muxVars: func(r *http.Request) map[string]string {
			// Using RequestURI to mock the different outputs
			switch r.RequestURI {
			case "notExist":
				return map[string]string{
					"institution": "inst1",
					"collection":  "coll1",
					"artwork":     "notExist",
				}
			}
			return map[string]string{
				"institution": "inst1",
				"collection":  "coll1",
				"artwork":     "pic1",
			}
		},
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name      string
		app       *artsPartsApp
		args      args
		expHeader int
	}{
		{
			"not exist post",
			app,
			args{
				httptest.NewRecorder(),
				&http.Request{
					RequestURI: "notExist",
					Method:     "GET",
				},
			},
			404,
		},
		{
			"exist post",
			app,
			args{
				httptest.NewRecorder(),
				&http.Request{
					RequestURI: "exist",
					Method:     "GET",
				},
			},
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.app.artwork(tt.args.w, tt.args.r)
			if tt.expHeader != tt.args.w.Result().StatusCode {
				t.Error("Wrong return code")
			}
		})
	}
}
