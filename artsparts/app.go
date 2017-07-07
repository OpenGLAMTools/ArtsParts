package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ChimeraCoder/anaconda"
	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/OpenGLAMTools/ArtsParts/helpers"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

// ArtsPartsApp contains all the handlefuncs
type ArtsPartsApp struct {
	artsparts        *artsparts.App
	conf             Conf
	muxVars          func(r *http.Request) map[string]string
	getSessionValues func(r *http.Request) map[string]string
}

// NewArtsPartsApp creates a new app
func NewArtsPartsApp(conf Conf) (*ArtsPartsApp, error) {
	apApp, err := artsparts.NewApp(conf.SourceFolder)
	if err != nil {
		return nil, err
	}
	app := &ArtsPartsApp{
		artsparts:        apApp,
		conf:             conf,
		muxVars:          mux.Vars,
		getSessionValues: getSessionValues,
	}
	return app, nil
}

func (app *ArtsPartsApp) defaultTemplateData(r *http.Request) *TemplateData {
	values := app.getSessionValues(r)

	admInst := app.artsparts.AdminInstitutions(values["twitter"])
	isAdmin := false
	if len(admInst) > 0 {
		isAdmin = true
	}
	vars := app.muxVars(r)
	if vars == nil {
		vars = make(map[string]string)
	}
	return &TemplateData{
		JSFiles:  []string{"/lib/app.js"},
		CSSFiles: []string{"/lib/custom.css"},
		JQuery:   true,
		VueJS:    false,
		Title:    "artsparts",
		User:     values["twitter"],
		Vars:     vars,
		Admin:    isAdmin,
		Session:  values,
	}
}
func (app *ArtsPartsApp) defaultFuncMap() template.FuncMap {
	funcMap := make(template.FuncMap)
	funcMap["vue"] = func(s string) string { return fmt.Sprintf("{{%s}}", s) }
	formatTS := func(ts, layout string) (string, error) {
		t, err := time.Parse(artsparts.TimneStampLayout, ts)
		if err != nil {
			return "", err
		}
		return t.Format(layout), nil
	}
	funcMap["formatTS"] = formatTS
	funcMap["tsToDateTime"] = func(s string) string {
		dt, _ := formatTS(s, "02.01.2006 15:04")
		return dt
	}
	return funcMap
}

func (app *ArtsPartsApp) executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	funcMap := app.defaultFuncMap()
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/*.tmpl.htm")
	if err != nil {
		log.Error("app.executeTemplate: error when parse glob: ", err)
	}
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		w.WriteHeader(404)
		w.Write([]byte("<h1>404 file not found</h1>"))
		log.Error("app.executeTemplate: error when execute template: ", err)
	}
}

// Page serves the templates direct. It is important to add a new template also
// to the allowed pages variable.
func (app *ArtsPartsApp) Page(w http.ResponseWriter, r *http.Request) {
	// Config the allowed pages here
	allowedPages := []string{"admin"}

	data := app.defaultTemplateData(r)
	page := data.Vars["page"]
	if !helpers.StringInSlice(page, allowedPages) {
		w.WriteHeader(404)
		w.Write([]byte("<h1>404 file not found</h1>"))
		return
	}

	// page individual configuration
	switch page {
	case "admin":
		data.AddJS("/lib/admin.js")
		data.VueJS = true
	}
	app.executeTemplate(w, page, data)
}

// Timeline serves the homepage with timeline
func (app *ArtsPartsApp) Timeline(w http.ResponseWriter, r *http.Request) {
	data := app.defaultTemplateData(r)
	var err error
	data.Timeline, err = app.artsparts.GetTimeline("")
	if err != nil {
		log.Error("app.timeline: error requesting timeline", err)
	}
	app.executeTemplate(w, "timeline", data)
}

// TweetNewArtworks checks if there is an artwork which is new inside the timeline
func (app *ArtsPartsApp) TweetNewArtworks() {
	ticker := time.Tick(time.Second * 30)
	for now := range ticker {
		tl, err := app.artsparts.GetTimeline("")
		if err != nil {
			log.Error("TweetNewArtwork: ", err)
		}
		for _, artw := range tl {
			if artw.TweetID == 0 {
				log.Info(artw)
				go app.TweetArtwork(artw)
			}
		}

		log.Info("Tick:", now)
	}
}

// TweetArtwork tweets an artwork
func (app *ArtsPartsApp) TweetArtwork(artw *artsparts.Artwork) {
	img, err := artw.Image()
	if err != nil {
		log.Error("TweeArtwork: artw.Image()", err)
		return
	}
	twitterAPI := anaconda.NewTwitterApi(
		getenv("ACCESS_TOKEN"),
		getenv("ACCESS_TOKEN_SECRET"),
	)
	twitterID, _, err := tweetImage(
		fmt.Sprintf("%s %s Neues ArtPart Bild verfügbar. %s%s",
			artw.InstitutionTwitter(),
			artw.HashTag,
			app.conf.URL,
			artw.ShortLink,
		),
		img,
		twitterAPI)
	if err != nil {
		log.Error("TweetArtwork: tweetImage: ", err)
		return
	}
	artw.TweetID = twitterID
	err = artw.WriteData()
	if err != nil {
		log.Error("TweetArtwork: error writing data: ", err)
	}
}

// Img is the handler for serving the images. The url accepts also different
// sizes. If size is part of the url the image is resized.
//   * small 150x150
//   * medium 300x300
//   * big 600x600
//   * huge 800x800
//   * massive 960x960
func (app *ArtsPartsApp) Img(w http.ResponseWriter, r *http.Request) {
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
		return
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
		return
	}
	switch size {
	case "small":
		img = imaging.Fit(img, 150, 150, imaging.Lanczos)
	case "medium":
		img = imaging.Fit(img, 300, 300, imaging.Lanczos)
	case "big":
		img = imaging.Fit(img, 600, 600, imaging.Lanczos)
	case "huge":
		img = imaging.Fit(img, 800, 800, imaging.Lanczos)
	case "massive":
		img = imaging.Fit(img, 960, 960, imaging.Lanczos)
	}
	err = imaging.Encode(w, img, imaging.JPEG)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Can not encode image"))
		log.Error("can not encode image", err)
	}
}

// Collection is the REST api for serving the Collection via json
func (app *ArtsPartsApp) Collection(w http.ResponseWriter, r *http.Request) {
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

// Institution is the REST api for serving the institution via json
func (app *ArtsPartsApp) Institution(w http.ResponseWriter, r *http.Request) {
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

// Editor is the handlefunc to serve the editor
func (app *ArtsPartsApp) Editor(w http.ResponseWriter, r *http.Request) {
	// path:
	// /editor/{institution}/{collection}/{artwork}
	/*vars := app.muxVars(r)
	instID := vars["institution"]
	collID := vars["collection"]
	artwID := vars["artwork"]*/
	data := app.defaultTemplateData(r)
	data.AddCSS("https://cdnjs.cloudflare.com/ajax/libs/cropper/2.3.4/cropper.css")
	data.AddJS("https://cdnjs.cloudflare.com/ajax/libs/cropper/2.3.4/cropper.js")
	data.AddJS("/lib/editor.js")
	// enable vuejs here
	data.VueJS = true
	artw := app.artworkFromVars(data.Vars, w)
	tmplData := struct {
		*TemplateData
		Artwork *artsparts.Artwork
	}{
		data,
		artw,
	}
	app.executeTemplate(w, "editor", tmplData)

}

// Artpart serves the json api for tweeting a created artpart
func (app *ArtsPartsApp) Artpart(w http.ResponseWriter, r *http.Request) {
	data := app.defaultTemplateData(r)
	if data.User == "" {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("artpart: error reading from body", err)
		return
	}
	ap := &artsparts.Part{
		User: data.User,
	}
	err = json.Unmarshal(rbody, ap)
	if err != nil {
		log.Error("artpart: error unmarshaling body", err)
		return
	}
	artw := app.artworkFromVars(data.Vars, w)
	img, err := artw.Artpart(ap)
	if err != nil {
		log.Error("artpart: error creating artpart image", err)
		return
	}
	twitterAPI := anaconda.NewTwitterApi(
		data.Session["access_token"],
		data.Session["access_token_secret"],
	)
	err = postPartTweet(ap, img, twitterAPI)
	if err != nil {
		log.Error("artpart: error post tweet", err)
		return
	}
	artw.AddPart(ap)
	err = artw.WriteData()
	if err != nil {
		log.Error("artpart: error WriteData", err)
		return
	}
	//imaging.Save(img, "artpart.jpg")
}

func (app *ArtsPartsApp) ArtworkPage(w http.ResponseWriter, r *http.Request) {
	data := app.defaultTemplateData(r)
	//data.AddJS("https://platform.twitter.com/widgets.js")
	data.AddJS("/lib/artwork.js")
	artw := app.artworkFromVars(data.Vars, w)
	tmplData := struct {
		*TemplateData
		Artwork *artsparts.Artwork
	}{
		data,
		artw,
	}
	app.executeTemplate(w, "artwork", tmplData)
}

// Artwork is the REST api for the AdminInstitution app
func (app *ArtsPartsApp) ArtworkData(w http.ResponseWriter, r *http.Request) {
	// path:
	// /data/{institution}/{collection}/{artwork}
	vars := app.muxVars(r)

	artw := app.artworkFromVars(vars, w)
	switch r.Method {
	case "POST":
		session := app.getSessionValues(r)

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

func (app *ArtsPartsApp) artworkFromVars(vars map[string]string, w http.ResponseWriter) *artsparts.Artwork {
	instID := vars["institution"]
	collID := vars["collection"]
	artwID := vars["artwork"]
	artw, ok := app.artsparts.GetArtwork(instID, collID, artwID)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Artwork not found"))
	}
	return artw
}

// AdminInstitutions is the rest api for serving the insitutions where the user is
// admin
func (app *ArtsPartsApp) AdminInstitutions(w http.ResponseWriter, r *http.Request) {
	session := app.getSessionValues(r)

	twitterName := session["twitter"]
	inss := app.artsparts.AdminInstitutions(twitterName)

	b, err := json.Marshal(inss)
	if err != nil {
		log.Error("error marshaling institution", err)
	}
	w.Write(b)
}
