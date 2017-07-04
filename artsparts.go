package artsparts

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/OpenGLAMTools/ArtsParts/helpers"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

// ConfFileName is the static value for all yaml conf files
const ConfFileName = "conf.yml"

// DataFileName defines the filename where data is stored.
const DataFileName = "data.json"

// TimeStampLayout defines the time layout string for the
// timestamp
const TimneStampLayout = "200601021504"

// Institutions holds the complete logic of the artsparts site.
// The insitutions are organized over a slice.
type Institutions []*Institution

// App wraps the whole logic for a simple communication from
// the http handlers.
type App struct {
	Institutions Institutions
}

// NewApp creates a new app by a given filepath
func NewApp(fpath string) (*App, error) {
	app := &App{}
	ls, err := ioutil.ReadDir(fpath)
	if err != nil {
		return app, errors.Errorf("%s\nPath:\n%s", err, fpath)
	}
	for _, d := range ls {
		if !d.IsDir() {
			continue
		}
		inst, err := NewInstitution(filepath.Join(fpath, d.Name()))
		if err != nil {
			return app, errors.Errorf("newInstitution error: %s\nFolder: %s",
				err,
				d.Name(),
			)
		}
		app.Institutions = append(app.Institutions, inst)
	}
	return app, nil
}

// GetInstitution returns the pointer to the given id. Second return
// parameter gives false, when no entry is found.
func (a *App) GetInstitution(instID string) (*Institution, bool) {
	for _, i := range a.Institutions {
		if i.ID == instID {
			return i, true
		}
	}
	return nil, false
}

// GetCollection returns the pointer to the given ids.
func (a *App) GetCollection(instID, collID string) (*Collection, bool) {
	i, ok := a.GetInstitution(instID)
	if !ok {
		return nil, false
	}
	c, ok := i.Collections[collID]
	return c, ok
}

// GetArtwork returns the pointer to the given ids.
func (a *App) GetArtwork(instID, collID, artwID string) (*Artwork, bool) {
	c, ok := a.GetCollection(instID, collID)
	if !ok {
		return nil, false
	}
	return c.GetArtwork(artwID)
}

// GetTimeline returns artworks, which can displayed in a timeline
// filter allows to get just special artowrks
// The logic is
// /[institution]/[collection]/[artwork.Name]
// The filter uses the simple regexp.MatchString() function
func (a *App) GetTimeline(filter string) ([]*Artwork, error) {
	tl := []*Artwork{}
	for _, inst := range a.Institutions {
		for _, coll := range inst.Collections {

			for _, artw := range coll.Artworks {
				p := fmt.Sprintf("/%s/%s/%s", inst.Name, coll.Name, artw.Name)

				match, err := regexp.MatchString(filter, p)
				if err != nil {
					return tl, errors.WithStack(err)
				}
				if match {
					if artw.Timestamp != "" {
						tl = append(tl, artw)
					}

				}
			}

		}
	}
	sort.Sort(Timeline(tl))
	return tl, nil
}

// AdminInstitutions returns all the intitutions, where the user is admin
func (a *App) AdminInstitutions(userName string) Institutions {
	ins := []*Institution{}
	for _, i := range a.Institutions {
		if helpers.StringInSlice(userName, i.Admins) {
			ins = append(ins, i)
		}
	}
	return ins
}

// Institution defines a museum or another partner of the site,
// which offers collections of digital art images.
//
// Following structure of the folders is needed:
//
// institution
// └───collections
//     └───collection1
//         ├───pic1
//         └───pic2
//
// The ID has to be unique and is always used inside the url.
type Institution struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	License     string                 `json:"license" yaml:"license"`
	Order       int                    `json:"order" yaml:"order"`
	Collections map[string]*Collection `json:"collections" yaml:"-"`
	Admins      []string               `json:"admins" yaml:"admins"`
}

func loadConf(confFile string, out interface{}) error {
	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, out)
	return nil
}

// NewInstitution takes a filepath and loads the configuration. Then
// it loads all the collections.
func NewInstitution(fpath string) (*Institution, error) {
	inst := &Institution{}
	confFile := filepath.Join(fpath, ConfFileName)
	if err := loadConf(confFile, inst); err != nil {
		return inst, err
	}
	if inst.ID == "" {
		inst.ID = filepath.Base(fpath)
	}
	inst.Collections = make(map[string]*Collection)
	ls, err := ioutil.ReadDir(fpath)
	if err != nil {
		return inst, err
	}
	for _, d := range ls {
		if !d.IsDir() {
			continue
		}
		coll, err := NewCollection(filepath.Join(fpath, d.Name()), inst)
		if err != nil {
			return inst, err
		}
		inst.Collections[coll.ID] = coll
	}
	return inst, nil
}

// Collection represents a group of artworks, which are presented
// together. It could be grouped after artist or a specific style.
// The Order property allows to sort the collections of a institution.
type Collection struct {
	ID          string     `json:"id" yaml:"id"`
	Name        string     `json:"name" yaml:"name"`
	Description string     `json:"description" yaml:"description"`
	License     string     `json:"license" yaml:"license"`
	Order       int        `json:"order" yaml:"order"`
	Artworks    []*Artwork `json:"artworks" yaml:"-"`
	institution *Institution
}

// NewCollection loads the configuration and creates a pointer to the
// new collection
func NewCollection(fpath string, inst *Institution) (*Collection, error) {
	coll := &Collection{
		institution: inst,
	}
	confFile := filepath.Join(fpath, ConfFileName)
	if err := loadConf(confFile, coll); err != nil {
		return coll, err
	}
	if coll.ID == "" {
		coll.ID = filepath.Base(fpath)
	}
	if coll.Name == "" {
		coll.Name = coll.ID
	}
	ls, err := ioutil.ReadDir(fpath)
	if err != nil {
		return coll, err
	}
	for _, d := range ls {
		if !d.IsDir() {
			continue
		}
		artw, err := NewArtwork(filepath.Join(fpath, d.Name()), coll)
		if err != nil {
			return coll, err
		}
		coll.Artworks = append(coll.Artworks, artw)
	}
	return coll, nil
}

// GetArtwork searches for an artwork id and returns a pointer to the
// artwork. If there is not artwork found with the given ID the second
// return value returns false.
func (coll *Collection) GetArtwork(artwID string) (*Artwork, bool) {
	for _, artw := range coll.Artworks {
		if artw.ID == artwID {
			return artw, true
		}
	}
	return nil, false
}

// Artwork is one element like for example a picture. The picture
// should be placed inside a folder. One artwork per folder! The
// foldername is then used as the id.
// The conf file here is stored as JSON, because the content is created
// and edited via a configuration dialog.
// Path contains the part of the url how the artwork can be found:
// /[institution]/[collection]/[artwork]
type Artwork struct {
	Timestamp   string `json:"timestamp"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URIPath     string `json:"-"`
	Parts       []*Part
	fpath       string
	collection  *Collection
}

// NewArtwork loads an artwork configuration and return a pointer.
func NewArtwork(fpath string, coll *Collection) (*Artwork, error) {
	artw := &Artwork{
		fpath:      fpath,
		collection: coll,
	}
	dataFilePath := filepath.Join(fpath, DataFileName)
	b, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return artw, err
		}
		// default values when data file is created
		artw.ID = filepath.Base(fpath)
		artw.Name = artw.ID
		// ensure the data file
		err = artw.WriteData()
		// return the fresh artwork
		return artw, err
	}
	if err := json.Unmarshal(b, artw); err != nil {
		return artw, err
	}
	artw.URIPath = fmt.Sprintf("/%s/%s/%s", coll.institution.ID, coll.ID, artw.ID)
	return artw, nil
}

// ImgFile return the filename of the image.
func (artw *Artwork) ImgFile() (string, error) {
	fileTypes := []string{".jpg", ".jpeg", ".png"}
	ls, err := ioutil.ReadDir(artw.fpath)
	if err != nil {
		return "", errors.WithMessage(err, "ReadDir in artwork.ImgFile()")
	}
	for _, f := range ls {
		if f.IsDir() {
			continue
		}
		if helpers.StringInSlice(filepath.Ext(f.Name()), fileTypes) {
			return f.Name(), nil
		}
	}
	return "", errors.New("No image file found")
}

// Artpart creates the part of the artwork. Every number is relative to the size
// of the picture. To get the x value in pixel you need to multiply it to the width
func (artw *Artwork) Artpart(x, y, width, height float64) (image.Image, error) {
	f, err := artw.ImgFile()
	if err != nil {
		return nil, err
	}
	img, err := imaging.Open(f)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	imgWidth := float64(bounds.Max.X)
	imgHeigth := float64(bounds.Max.Y)
	rect := image.Rect(
		int(x*imgWidth),
		int(y*imgHeigth),
		int((x+width)*imgWidth),
		int((y+height)*imgHeigth))
	artp := imaging.Crop(img, rect)
	return artp, nil
}

// WriteData writes the artw into a file. If the
// file does not exist the file is created.
func (artw *Artwork) WriteData() error {
	b, err := artw.Marshal()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(artw.dataFilePath(), b, 0777)
}

// Marshal wraps the json marshal func. If the file format
// should be changed it can be done here.
func (artw *Artwork) Marshal() ([]byte, error) {
	return json.MarshalIndent(artw, "", "   ")
}

// Path returns the stored path to the artwork as string
func (artw *Artwork) Path() string {
	return artw.fpath
}

// IsAdminUser returns true, when the given user has admin rights
func (artw *Artwork) IsAdminUser(userName string) bool {
	return helpers.StringInSlice(userName, artw.collection.institution.Admins)
}

func (artw *Artwork) dataFilePath() string {
	return filepath.Join(artw.fpath, DataFileName)
}

// Part represends a part, which is tweeted from artsparts
type Part struct {
	TopLeft     int
	BottomRight int
	TweetID     int64
	MediaID     int64
}

type Timeline []*Artwork

func (tl Timeline) Len() int           { return len(tl) }
func (tl Timeline) Swap(i, j int)      { tl[i], tl[j] = tl[j], tl[i] }
func (tl Timeline) Less(i, j int) bool { return tl[i].Timestamp > tl[j].Timestamp }

type TimelineItem struct {
	InsitutionName   string
	CollectionName   string
	ArtworkName      string
	ArtworkTimestamp int
}

// User of the artsparts page. Everything about the user goes over
// twitter. A twitter account and the authentication to the account
// is compulsary to use artsparts.
// The TwitterID is a unique number, which is normaly not known by
// the user. The TwitterName is the screen_name
type User struct {
	TwitterID   int64
	TwitterName string
	Email       string
}
