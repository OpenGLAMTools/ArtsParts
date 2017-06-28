package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type artsPartsApp struct {
	artsparts *artsparts.App
}

func newArtsPartsApp(fpath string) (*artsPartsApp, error) {
	apApp, err := artsparts.NewApp(fpath)
	if err != nil {
		return nil, err
	}
	app := &artsPartsApp{apApp}
	return app, nil
}

func (app *artsPartsApp) artwork(w http.ResponseWriter, r *http.Request) {
	// path:
	// /data/{institution}/{collection}/{artwork}
	vars := mux.Vars(r)
	instID := vars["institution"]
	collID := vars["collection"]
	artwID := vars["artwork"]
	artw, ok := app.artsparts.GetArtwork(instID, collID, artwID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Artwork not found"))
	}
	b, err := artw.Marshal()
	if err != nil {
		log.Error("error marshaling artwork", err)
	}
	w.Write(b)
	session, _ := gothic.Store.Get(r, sessionName)
	fmt.Fprintf(w, "%#v", session.Values)
}

func (app *artsPartsApp) img(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instID := vars["institution"]
	collID := vars["collection"]
	artwID := vars["artwork"]
	artw, ok := app.artsparts.GetArtwork(instID, collID, artwID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Artwork not found"))
	}
	imgFile, err := artw.ImgFile()
	if err != nil {
		log.Error("app.img() artw.ImgFile: ", err)
	}
	b, err := ioutil.ReadFile(filepath.Join(artw.Fpath, imgFile))
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Can not load image"))
		log.Error("can not load image file", err)
	} else {
		w.Write(b)
	}

}

func (app *artsPartsApp) collection(w http.ResponseWriter, r *http.Request) {
	// path:
	// /data/{institution}/{collection}/{artwork}
	vars := mux.Vars(r)
	instID := vars["institution"]
	collID := vars["collection"]
	coll, ok := app.artsparts.GetCollection(instID, collID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Collection not found"))
	}
	b, err := json.Marshal(coll)
	if err != nil {
		log.Error("error marshaling artwork", err)
	}
	w.Write(b)
}

func (app *artsPartsApp) institution(w http.ResponseWriter, r *http.Request) {
	// path:
	// /data/{institution}/{collection}/{artwork}
	vars := mux.Vars(r)
	instID := vars["institution"]
	inst, ok := app.artsparts.GetInstitution(instID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Institution not found"))
	}
	b, err := json.Marshal(inst)
	if err != nil {
		log.Error("error marshaling institution", err)
	}
	w.Write(b)
}

func (app *artsPartsApp) adminInstitutions(w http.ResponseWriter, r *http.Request) {
	session, err := getSessionValues(r)
	if err != nil {
		log.Error("adminInstitutions error getSessionValues:", err)
	}
	twitterName := session["twitter"]
	inss := app.artsparts.AdminInstitutions(twitterName)

	b, err := json.Marshal(inss)
	if err != nil {
		log.Error("error marshaling institution", err)
	}
	w.Write(b)
}
