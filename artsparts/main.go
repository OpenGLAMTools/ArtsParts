package main

import (
	"fmt"
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
	initTwitter()
	r := makeRoutes()
	r = addAuthRoutes(r)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is artsParts")
	})
	r.HandleFunc("/tweet", postTweetHandler)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}

func makeRoutes() *mux.Router {
	return mux.NewRouter()
}
