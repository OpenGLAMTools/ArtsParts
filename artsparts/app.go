package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/OpenGLAMTools/ArtsParts/helpers"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

// pagesFileName defines the file, where the static pages are defined.
const pagesFileName = "pages.md"

// ArtsPartsApp contains all the handlefuncs
type ArtsPartsApp struct {
	artsparts        *artsparts.App
	conf             Conf
	muxVars          func(r *http.Request) map[string]string
	getSessionValues func(r *http.Request) map[string]string
	env              map[string]string
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
		env:              conf.Env,
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
	pages, err := loadPages(pagesFileName)
	if err != nil {
		log.Error("defaultTemplateData: error loading pages ", err)
	}

	return &TemplateData{
		JSFiles:         []string{"/lib/app.js"},
		CSSFiles:        []string{"/lib/custom.css"},
		GoogleAnalytics: app.env["GOOGLE_ANALYTICS"],
		JQuery:          true,
		VueJS:           false,
		Title:           "",
		User:            values["twitter"],
		Vars:            vars,
		Pages:           pages,
		Artsparts:       app.artsparts,
		Admin:           isAdmin,
		Session:         values,
	}
}
func (app *ArtsPartsApp) defaultFuncMap() template.FuncMap {
	funcMap := make(template.FuncMap)
	funcMap["vue"] = func(s string) string { return fmt.Sprintf("{{%s}}", s) }
	formatTS := func(ts, layout string) (string, error) {
		t, err := time.Parse(artsparts.TimeStampLayout, ts)
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
	funcMap["md"] = func(s string) template.HTML {
		return template.HTML(blackfriday.MarkdownBasic([]byte(s)))
	}
	funcMap["inPage"] = func(index, pagenr int) bool {
		min := (pagenr - 1) * app.conf.ItemsPerPage
		max := min + app.conf.ItemsPerPage
		return min <= index && index < max
	}
	funcMap["pageExists"] = func(pagenr, items int) bool {
		min := (pagenr - 1) * app.conf.ItemsPerPage
		return min < items
	}
	funcMap["add"] = func(a, b int) int {
		return a + b
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
	pages, err := loadPages(pagesFileName)
	if err != nil {
		log.Error("app.Page() error loading pages: ", err)
	}
	// Config the allowed pages here
	allowedPages := []string{"admin"}
	// add the pages from the conf
	for _, p := range pages {
		allowedPages = append(allowedPages, p.Path)
	}

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
	default:
		thisPage := getPage(page, pages)
		data.Title = thisPage.Title
		data.Vars["title"] = thisPage.Title
		data.Vars["text"] = thisPage.Text
		page = "page"
	}
	app.executeTemplate(w, page, data)
}

// Timeline serves the homepage with timeline
func (app *ArtsPartsApp) Timeline(w http.ResponseWriter, r *http.Request) {
	data := app.defaultTemplateData(r)
	data.Title = "Timeline alle Sammlungen"
	filter := ""
	institution, ok := data.Vars["institution"]
	if ok {
		filter = fmt.Sprintf("/%s/*", institution)
		inst, ok := app.artsparts.GetInstitution(institution)
		if !ok {
			w.WriteHeader(404)
			return
		}
		data.Title = fmt.Sprintf("Timeline - %s", inst.Name)
		collection, ok := data.Vars["collection"]
		if ok {
			coll, ok := app.artsparts.GetCollection(institution, collection)
			if !ok {
				w.WriteHeader(404)
				return
			}
			filter = fmt.Sprintf("/%s/%s/*", institution, collection)
			data.Title = fmt.Sprintf("Timeline - %s by %s", coll.Name, inst.Name)
		}
	}
	q := r.URL.Query()
	pagenr, err := strconv.Atoi(q.Get("page"))
	// If value can not be converted set page 1 as default
	if err != nil {
		pagenr = 1
	}
	data.Pagenr = pagenr
	data.Timeline, err = app.artsparts.GetPublishedTimeline(filter)

	if err != nil {
		log.Error("app.timeline: error requesting timeline", err)
	}
	app.executeTemplate(w, "timeline", data)
}

// TweetNewArtworks checks if there is an artwork which is new inside the timeline
func (app *ArtsPartsApp) TweetNewArtworks() {
	ticker := time.Tick(time.Minute)
	for _ = range ticker {
		tl, err := app.artsparts.GetPublishedTimeline("")
		if err != nil {
			log.Error("TweetNewArtwork: ", err)
		}
		for _, artw := range tl {
			if artw.TweetID == 0 {
				go app.TweetArtwork(artw)
			}
		}
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
		app.env["ACCESS_TOKEN"],
		app.env["ACCESS_TOKEN_SECRET"],
	)
	tweetResponse, err := tweetImage(
		fmt.Sprintf("Neues ArtPart Bild verfügbar. %s%s by %s #ArtsParts %s",
			app.conf.URL,
			artw.ShortLink,
			artw.InstitutionTwitter(),
			artw.HashTag,
		),
		img,
		twitterAPI)
	if err != nil {
		log.Error("TweetArtwork: tweetImage: ", err)
		return
	}
	log.Infoln("New Artwork. TweetID", tweetResponse.twitterIDString)
	artw.TweetID = tweetResponse.twitterID
	artw.TweetIDString = tweetResponse.twitterIDString
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
	fpath := filepath.Join(artw.Path(), imgFile)
	img, err := getImageWithSize(fpath, size)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Can not load image"))
		log.Error("can not load image file", err)
		return
	}
	err = imaging.Encode(w, img, imaging.JPEG)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Can not encode image"))
		log.Error("can not encode image", err)
	}
}

func getImageWithSize(fpath, size string) (image.Image, error) {
	cachePath := fmt.Sprintf("%s-%s", fpath, size)
	sum := sha256.Sum256([]byte(cachePath))
	cacheString := fmt.Sprintf("%x", sum)
	cacheFile := filepath.Join(
		"cache",
		cacheString[:2],
		cacheString,
	)
	if helpers.FileExists(cacheFile) {
		return imaging.Open(cacheFile)
	}

	img, err := imaging.Open(fpath)
	if err != nil {
		return nil, err
	}
	switch size {
	case "mini":
		img = imaging.Fit(img, 35, 35, imaging.Lanczos)
	case "tiny":
		img = imaging.Fit(img, 80, 80, imaging.Lanczos)
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
	// Write the fitted image into the cache
	err = os.MkdirAll(
		filepath.Dir(cacheFile),
		0777,
	)
	if err != nil {
		log.Error("error creating cache dir: ", err)
	}
	cf, err := os.Create(cacheFile)
	defer cf.Close()
	if err != nil {
		log.Error("error creating cache file:", err)
	}
	err = imaging.Encode(cf, img, imaging.JPEG)
	if err != nil {
		log.Error("error writing cache file", err)
	}
	return img, nil
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
	data.Title = fmt.Sprintf("Editor: %s", artw.Name)
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
	data.Title = artw.Name
	tmplData := struct {
		*TemplateData
		Artwork *artsparts.Artwork
	}{
		data,
		artw,
	}
	app.executeTemplate(w, "artwork", tmplData)
}

// ArtworkData is the REST api for the AdminInstitution app
// Security isssue when a POST request is send.
// Not all data is allowed to be changed over the POST request. Just the data
// which the admin can change
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
		newData := &artsparts.Artwork{}
		err = json.Unmarshal(rbody, newData)
		if err != nil {
			log.Error("artwork: error unmarshaling body", err)
			return
		}
		artsparts.CopyArtwork(artw, newData, "admin")
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
	// reverse the order
	if artw != nil && len(artw.Parts) > 0 {
		var parts []*artsparts.Part
		for i := range artw.Parts {
			parts = append(parts, artw.Parts[len(artw.Parts)-i-1])
		}
		artw.Parts = parts
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
