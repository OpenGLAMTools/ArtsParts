package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var confFile = ".conf.yml"

var log = logrus.New()

func main() {
	conf, err := loadConf(confFile)
	if err != nil {
		log.Fatal("Can not load conf: ", confFile, err)
	}
	log.Level = logrus.Level(conf.LogLevel)
	initAuth(conf)
	initTwitter()
	r := makeRoutes()
	r = addAuthRoutes(r)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is artsParts")
	})
	r.HandleFunc("/tweet", postTweetHandler)
	log.Infoln("Starting server at: ", conf.ServerPort)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}

func makeRoutes() *mux.Router {
	return mux.NewRouter()
}
