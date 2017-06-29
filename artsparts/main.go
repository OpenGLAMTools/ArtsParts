package main

import (
	"bytes"
	"html/template"
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
	r.HandleFunc("/page/{page}", pageHandler)
	r.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("templates/lib"))))
	r.HandleFunc("/tweet", postTweetHandler)
	// Auth routes
	// /auth/twitter
	r = addAuthRoutes(r)

	app, err := newArtsPartsApp(conf.SourceFolder)
	if err != nil {
		log.Fatal("error initializing app:", err)
	}
	r.HandleFunc("/", app.timeline)
	r.HandleFunc("/data/admin", app.adminInstitutions).Methods("GET")
	r.HandleFunc("/data/{institution}/{collection}/{artwork}", app.artwork).Methods("GET", "POST")
	r.HandleFunc("/img/{institution}/{collection}/{artwork}", app.img).Methods("GET")
	r.HandleFunc("/data/{institution}/{collection}", app.collection).Methods("GET")
	r.HandleFunc("/data/{institution}", app.institution).Methods("GET")

	log.Infoln("Starting server at: ", conf.ServerPort)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	tmpl, err := template.ParseGlob("templates/*.tmpl.htm")
	if err != nil {
		log.Error("error when parse glob: ", err)
	}
	values, err := getSessionValues(r)
	if err != nil {
		log.Warningln("pageHandler: Error when getSessionValues:", err)
	}
	data := templateData{
		JSFiles:  []string{"app.js", "admin.js"},
		CSSFiles: []string{"custom.css"},
		JQuery:   true,
		VueJS:    true,
		User:     values["twitter"],
	}
	if err := tmpl.ExecuteTemplate(w, page, data); err != nil {
		w.WriteHeader(404)
		w.Write(bytes.NewBufferString("<h1>404 file not found</h1>").Bytes())
		log.Info("error when execute template: ", err)
	}
}
