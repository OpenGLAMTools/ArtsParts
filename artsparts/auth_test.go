package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestAddAuthRoutes(t *testing.T) {
	router := mux.NewRouter()
	router = addAuthRoutes(router)
	ts := httptest.NewServer(router)
	testUrls := []string{
		"/auth",
		"/auth/callback",
	}
	for _, tURL := range testUrls {
		res, err := http.Get(ts.URL + tURL)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != 200 {
			t.Errorf("Did not find %s", tURL)
		}
	}

}
