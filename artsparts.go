package artsparts

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/OpenGLAMTools/ArtsParts/helpers"
	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

// ConfFileName is the static value for all yaml conf files
const ConfFileName = "conf.yml"

// DataFileName defines the filename where data is stored.
const DataFileName = "data.json"

// TimeStampLayout defines the time layout string for the
// timestamp
//const TimeStampLayout = "200601021504 -0700 MST"
const TimeStampLayout = "200601021504"

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

// GetPublishedTimeline calls GetTimeline() and filters all artworks which
// has a timestamp in the future. So just the published artworks are returned
func (a *App) GetPublishedTimeline(filter string) ([]*Artwork, error) {
	var publishedTimeline []*Artwork
	tl, err := a.GetTimeline(filter)
	if err != nil {
		return publishedTimeline, err
	}
	now := time.Now()
	loc, err := time.LoadLocation("Local")
	for _, artw := range tl {
		artwTime, err := time.ParseInLocation(TimeStampLayout, artw.Timestamp, loc)
		if err != nil {
			return publishedTimeline, errors.WithStack(err)
		}
		if now.After(artwTime) {
			publishedTimeline = append(publishedTimeline, artw)
		}
	}
	return publishedTimeline, nil
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
	TwitterName string                 `json:"twittername" yaml:"twitter_name"`
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

// Part represends a part, which is tweeted from artsparts all coordinates
// are relative to get the pixel the values have to be multiplied with
// the image size. For example X * ImageWidth = Pixel for X
type Part struct {
	Text    string  `json:"tweettext"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	User    string  `json:"user"`
	TweetID int64   `json:"tweet_id,omitempty"`
	MediaID int64   `json:"media_id,omitempty"`
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
