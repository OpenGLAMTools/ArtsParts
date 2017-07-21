package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/OpenGLAMTools/ArtsParts/shortlink"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var log = logrus.New()

func main() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal("Can not get executable: ", err)
	}
	defaultConfPath := filepath.Join(
		filepath.Dir(executable),
		".conf.yml",
	)
	confFile := flag.String("conf", defaultConfPath, "Path to the configuration")
	flag.Parse()
	conf, err := loadConf(*confFile)
	if err != nil {
		log.Fatal("Can not load conf: ", confFile, err)
	}
	log.Level = logrus.Level(conf.LogLevel)
	initAuth(conf)
	initTwitter(conf)

	r := mux.NewRouter()

	r.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("templates/lib"))))

	// Auth routes
	// /auth/twitter
	r = shortlink.AddRoute(r)
	r = addAuthRoutes(r)

	app, err := NewArtsPartsApp(conf)
	if err != nil {
		log.Fatal("error initializing app:", err)
	}
	go app.TweetNewArtworks()
	r = addAppRoutes(r, app)

	log.Infoln("Starting server at: ", conf.ServerPort)
	log.Fatal(http.ListenAndServe(conf.ServerPort, r))

}

func addAppRoutes(r *mux.Router, app *ArtsPartsApp) *mux.Router {

	r.HandleFunc("/", app.Timeline)
	r.HandleFunc("/page/{page}", app.Page)
	r.HandleFunc("/data/admin", app.AdminInstitutions).Methods("GET")
	r.HandleFunc("/data/{institution}/{collection}/{artwork}", app.ArtworkData).Methods("GET", "POST")
	r.HandleFunc("/img/{institution}/{collection}/{artwork}", app.Img).Methods("GET")
	r.HandleFunc("/data/{institution}/{collection}", app.Collection).Methods("GET")
	r.HandleFunc("/data/{institution}", app.Institution).Methods("GET")

	r.HandleFunc("/editor/{institution}/{collection}/{artwork}", app.Editor).Methods("GET")
	r.HandleFunc("/artpart/{institution}/{collection}/{artwork}", app.Artpart).Methods("POST")
	r.HandleFunc("/artwork/{institution}/{collection}/{artwork}", app.ArtworkPage).Methods("GET")
	return r
}
