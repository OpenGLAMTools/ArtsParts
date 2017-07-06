package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAppRoutes(t *testing.T) {
	conf, _ := loadConf("default.conf.yml")
	router := mux.NewRouter()
	router = addAppRoutes(router, conf)
	ts := httptest.NewServer(router)
	tests := []struct {
		route         string
		expStatusCode int
		contains      string
	}{
		{"/", 200, "Timeline"},
		{"/page/admin", 200, "admin.js"},
		{"/page/nothing", 404, "404"},
		{"/data/inst1", 200, "inst1"},
		{"/data/inst123", 404, ""},
		{"/data/inst1/coll1", 200, "Modern Art"},
		{"/data/inst123/coll1", 404, ""},
		{"/data/inst1/coll222", 404, ""},
		{"/data/inst1/coll1/pic1", 200, "First Picture"},
		{"/data/inst1/coll1/pic123", 404, ""},
		{"/data/inst1/coll123/pic1", 404, ""},
		{"/img/inst1/coll1/pic1", 200, ""},
		{"/img/inst1/coll1/pic1?size=small", 200, ""},
		{"/img/inst1/coll1/pic1?size=medium", 200, ""},
		{"/img/inst1/coll1/pic1?size=big", 200, ""},
		{"/img/inst1/coll1/pic1?size=huge", 200, ""},
		{"/img/inst1/coll1/pic1?size=massive", 200, ""},
		{"/img/inst1/coll123/pic1", 404, ""},
		{"/img/inst1/coll1/pic123", 404, ""},
		{"/editor/inst1/coll1/pic1", 200, "cropper.js"},
		{"/editor/inst1/coll1/pic1", 200, "https://unpkg.com/vue"},
	}

	for _, tt := range tests {
		res, err := http.Get(ts.URL + tt.route)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != tt.expStatusCode {
			t.Errorf("Wrong Status code for %s\nGot: %03d\nExp: %03d\n",
				tt.route,
				res.StatusCode,
				tt.expStatusCode)
		}
		body, _ := ioutil.ReadAll(res.Body)
		assert.Contains(t, string(body), tt.contains)
	}

}
