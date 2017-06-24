package main

import (
	"bytes"
	"fmt"
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
	r := makeRoutes()
	r = addAuthRoutes(r)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is artsParts")
	})
	r.HandleFunc("/page/{page}", pageHandler)
	r.PathPrefix(`/lib/`).
		Handler(
			http.StripPrefix("/lib/",
				http.FileServer(http.Dir("templates/lib"))))
	r.HandleFunc("/tweet", postTweetHandler)
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
	data := templateData{
		JSFiles:  []string{"/app.js", "admin.js"},
		CSSFiles: []string{"custom.css"},
		JQuery:   true,
		VueJS:    true,
	}
	if err := tmpl.ExecuteTemplate(w, page, data); err != nil {
		w.WriteHeader(404)
		w.Write(bytes.NewBufferString("<h1>404 file not found</h1>").Bytes())
		log.Info("error when execute template: ", err)
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {

}

func makeRoutes() *mux.Router {
	return mux.NewRouter()
}
