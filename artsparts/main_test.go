package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestPagesRoute(t *testing.T) {

	router := mux.NewRouter()
	router = addAppRoutes(router, "../test")
	ts := httptest.NewServer(router)
	tests := []struct {
		route         string
		expStatusCode int
		contains      string
	}{
		{"/", 200, "Timeline"},
		{"/page/admin", 200, "admin.js"},
		{"/page/nothing", 404, "404"},
	}

	for _, tt := range tests {
		res, err := http.Get(ts.URL + tt.route)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != tt.expStatusCode {
			t.Errorf("Did not find %s\n", tt.route)
		}
		body, _ := ioutil.ReadAll(res.Body)
		assert.Contains(t, string(body), tt.contains)
	}

}
