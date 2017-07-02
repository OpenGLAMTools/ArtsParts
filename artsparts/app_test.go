package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

func Test_artsPartsApp_defaultTemplateData(t *testing.T) {
	app, _ := NewArtsPartsApp("../test")
	app.getSessionValues = func(r *http.Request) (map[string]string, error) {
		return map[string]string{
			"twitter": "user1",
		}, nil
	}
	app2, _ := NewArtsPartsApp("../test")
	app2.getSessionValues = func(r *http.Request) (map[string]string, error) {
		return map[string]string{
			"twitter": "user11",
		}, nil
	}

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		app  *ArtsPartsApp
		args args
		want templateData
	}{
		{
			"Admin user",
			app,
			args{nil},
			templateData{
				JSFiles:  []string{"app.js"},
				CSSFiles: []string{"custom.css"},
				JQuery:   true,
				VueJS:    false,
				Title:    "artsparts",
				User:     "user1",
				Admin:    true,
			},
		},
		{
			"Normal user",
			app2,
			args{nil},
			templateData{
				JSFiles:  []string{"app.js"},
				CSSFiles: []string{"custom.css"},
				JQuery:   true,
				VueJS:    false,
				Title:    "artsparts",
				User:     "user11",
				Admin:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.app.defaultTemplateData(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("artsPartsApp.defaultTemplateData() = \n%#v, want \n%#v", got, tt.want)
			}
		})
	}
}

func Test_artsPartsApp_artwork(t *testing.T) {
	// TODO: test is not complete
	artspartsApp, _ := artsparts.NewApp("../test")
	app := &ArtsPartsApp{
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
		app       *ArtsPartsApp
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
			tt.app.Artwork(tt.args.w, tt.args.r)
			if tt.expHeader != tt.args.w.Result().StatusCode {
				t.Error("Wrong return code")
			}
		})
	}
}
