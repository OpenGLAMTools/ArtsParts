package main

import (
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

	r := mux.NewRouter()

	r.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("templates/lib"))))
	r.HandleFunc("/tweet", postTweetHandler)
	// Auth routes
	// /auth/twitter
	r = addAuthRoutes(r)

	app, err := NewArtsPartsApp(conf.SourceFolder)
	if err != nil {
		log.Fatal("error initializing app:", err)
	}

	r.HandleFunc("/", app.Timeline)
	r.HandleFunc("/page/{page}", app.Page)
	r.HandleFunc("/data/admin", app.AdminInstitutions).Methods("GET")
	r.HandleFunc("/data/{institution}/{collection}/{artwork}", app.Artwork).Methods("GET", "POST")
	r.HandleFunc("/img/{institution}/{collection}/{artwork}", app.Img).Methods("GET")
	r.HandleFunc("/data/{institution}/{collection}", app.Collection).Methods("GET")
	r.HandleFunc("/data/{institution}", app.Institution).Methods("GET")

	r.HandleFunc("/editor/{institution}/{collection}/{artwork}", app.Editor).Methods("GET")
	r.HandleFunc("/artpart/{institution}/{collection}/{artwork}", app.Artpart).Methods("POST")

	log.Infoln("Starting server at: ", conf.ServerPort)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}
