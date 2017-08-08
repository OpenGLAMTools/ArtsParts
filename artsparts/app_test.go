package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/stretchr/testify/assert"
)

func Test_artsPartsApp_defaultTemplateData(t *testing.T) {
	conf, _ := loadConf("default.conf.yml")
	app, _ := NewArtsPartsApp(conf)
	app.getSessionValues = func(r *http.Request) map[string]string {
		return map[string]string{
			"twitter": "user1",
		}
	}
	app2, _ := NewArtsPartsApp(conf)
	app2.getSessionValues = func(r *http.Request) map[string]string {
		return map[string]string{
			"twitter": "user11",
		}
	}
	pages, _ := loadPages(pagesFileName)
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		app  *ArtsPartsApp
		args args
		want TemplateData
	}{
		{
			"Admin user",
			app,
			args{&http.Request{}},
			TemplateData{
				JSFiles:   []string{"/lib/app.js"},
				CSSFiles:  []string{"/lib/custom.css"},
				JQuery:    true,
				VueJS:     false,
				Title:     "",
				User:      "user1",
				Vars:      map[string]string{},
				Pages:     pages,
				Artsparts: app.artsparts,
				Admin:     true,
				Session: map[string]string{
					"twitter": "user1",
				},
			},
		},
		{
			"Normal user",
			app2,
			args{&http.Request{}},
			TemplateData{
				JSFiles:   []string{"/lib/app.js"},
				CSSFiles:  []string{"/lib/custom.css"},
				JQuery:    true,
				VueJS:     false,
				Title:     "",
				User:      "user11",
				Vars:      map[string]string{},
				Pages:     pages,
				Artsparts: app2.artsparts,
				Admin:     false,
				Session: map[string]string{
					"twitter": "user11",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.app.defaultTemplateData(tt.args.r); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("artsPartsApp.defaultTemplateData() = \n%#v, want \n%#v", *got, tt.want)
			}
		})
	}
}

func XXXTest_artsPartsApp_artwork(t *testing.T) {
	// TODO: test is not complete
	artspartsApp, _ := artsparts.NewApp("../test")
	app := &ArtsPartsApp{
		artsparts: artspartsApp,
		getSessionValues: func(r *http.Request) map[string]string {
			return map[string]string{
				"twitter": "user1",
			}
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
			//tt.app.Artwork(tt.args.w, tt.args.r)
			if tt.expHeader != tt.args.w.Result().StatusCode {
				t.Error("Wrong return code")
			}
		})
	}
}

func TestDafultFuncMap(t *testing.T) {
	conf, _ := loadConf("default.conf.yml")
	app, _ := NewArtsPartsApp(conf)
	fm := app.defaultFuncMap()
	vue := fm["vue"].(func(string) string)("abc")
	exp := "{{abc}}"
	if vue != exp {
		t.Errorf("Got: %s\nExp: %s\n", vue, exp)
	}
	tests := []struct {
		ts        string
		layout    string
		expError  bool
		expString string
	}{
		{
			"201301241155",
			"15:04 02.01.2006",
			false,
			"11:55 24.01.2013",
		},
		{
			"20130124115599",
			"15:04 02.01.2006",
			true,
			"",
		},
	}
	for _, tt := range tests {
		tstring, err := fm["formatTS"].(func(string, string) (string, error))(tt.ts, tt.layout)
		assert.Equal(t, tstring, tt.expString)

		assert.Equal(t, tt.expError, err != nil)
	}

}
