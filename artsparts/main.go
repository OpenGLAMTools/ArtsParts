package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var conf Conf

func init() {
	conf = loadConf()
}

func main() {
	r := makeRoutes()
	r = addAuthRoutes(r)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}

func makeRoutes() *mux.Router {
	return mux.NewRouter()
}
