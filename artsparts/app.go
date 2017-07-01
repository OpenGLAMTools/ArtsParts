package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

type artsPartsApp struct {
	artsparts        *artsparts.App
	muxVars          func(r *http.Request) map[string]string
	getSessionValues func(r *http.Request) (map[string]string, error)
}

func newArtsPartsApp(fpath string) (*artsPartsApp, error) {
	apApp, err := artsparts.NewApp(fpath)
	if err != nil {
		return nil, err
	}
	app := &artsPartsApp{
		apApp,
		mux.Vars,
		getSessionValues,
	}
	return app, nil
}

func (app *artsPartsApp) defaultTemplateData(r *http.Request) templateData {
	values, err := app.getSessionValues(r)
	if err != nil {
		log.Warningln("defaultTemplateData: Error when getSessionValues:", err)
	}
	return templateData{
		JSFiles:  []string{"app.js"},
		CSSFiles: []string{"custom.css"},
		JQuery:   true,
		VueJS:    false,
		User:     values["twitter"],
	}
}

func (app *artsPartsApp) executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseGlob("templates/*.tmpl.htm")
	if err != nil {
		log.Error("app.executeTemplate: error when parse glob: ", err)
	}
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		w.WriteHeader(404)
		w.Write([]byte("<h1>404 file not found</h1>"))
		log.Error("app.executeTemplate: error when execute template: ", err)
	}
}

func (app *artsPartsApp) timeline(w http.ResponseWriter, r *http.Request) {
	data := app.defaultTemplateData(r)
	var err error
	data.Timeline, err = app.artsparts.GetTimeline("")
	if err != nil {
		log.Error("app.timeline: error requesting timeline", err)
	}
	app.executeTemplate(w, "timeline", data)
}

func (app *artsPartsApp) artwork(w http.ResponseWriter, r *http.Request) {
	// path:
	// /data/{institution}/{collection}/{artwork}
	vars := app.muxVars(r)
	instID := vars["institution"]
	collID := vars["collection"]
	artwID := vars["artwork"]

	artw, ok := app.artsparts.GetArtwork(instID, collID, artwID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Artwork not found"))
	}
	switch r.Method {
	case "POST":
		session, err := app.getSessionValues(r)
		if err != nil {
			log.Error("artwork: error reading session", err)
			return
		}
		if !artw.IsAdminUser(session["twitter"]) {
			w.WriteHeader(403)
			w.Write([]byte("Forbidden"))
		}
		rbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("artwork: error reading from body", err)
			return
		}
		err = json.Unmarshal(rbody, artw)
		if err != nil {
			log.Error("artwork: error unmarshaling body", err)
			return
		}
		err = artw.WriteData()
		if err != nil {
			log.Error("artwork: error writing data", err)
			return
		}
	case "GET":
		b, err := artw.Marshal()
		if err != nil {
			log.Error("error marshaling artwork", err)
		}
		w.Write(b)
	}

}

// Img is the handler for serving the images. The url accepts also different
// sizes. If size is part of the url the image is resized.
//   * small 150x150
//   * medium 300x300
//   * big 600x600
func (app *artsPartsApp) img(w http.ResponseWriter, r *http.Request) {
	vars := app.muxVars(r)
	instID := vars["institution"]
	collID := vars["collection"]
	artwID := vars["artwork"]
	q := r.URL.Query()
	size := q.Get("size")
	artw, ok := app.artsparts.GetArtwork(instID, collID, artwID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Artwork not found"))
	}
	imgFile, err := artw.ImgFile()
	if err != nil {
		log.Error("app.img() artw.ImgFile: ", err)
	}
	img, err := imaging.Open(filepath.Join(artw.Path(), imgFile))
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Can not load image"))
		log.Error("can not load image file", err)
	}
	switch size {
	case "small":
		img = imaging.Fit(img, 150, 150, imaging.Lanczos)
	case "medium":
		img = imaging.Fit(img, 300, 300, imaging.Lanczos)
	case "big":
		img = imaging.Fit(img, 600, 600, imaging.Lanczos)
	}
	err = imaging.Encode(w, img, imaging.JPEG)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Can not encode image"))
		log.Error("can not encode image", err)
	}
}

func (app *artsPartsApp) collection(w http.ResponseWriter, r *http.Request) {
	// path:
	// /data/{institution}/{collection}/{artwork}
	vars := app.muxVars(r)
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
	vars := app.muxVars(r)
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
	session, err := app.getSessionValues(r)
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
