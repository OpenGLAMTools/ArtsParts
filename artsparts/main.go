package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var confFile = ".conf.yml"

func main() {
	conf, err := loadConf(confFile)
	if err != nil {
		log.Fatal("Can not load conf: ", confFile, err)
	}
	initAuth(conf)
	r := makeRoutes()
	r = addAuthRoutes(r)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}

func makeRoutes() *mux.Router {
	return mux.NewRouter()
}
