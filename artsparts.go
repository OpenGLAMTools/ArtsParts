package artsparts

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ConfFileName is the static value for all yaml conf files
const ConfFileName = "conf.yml"

// Institutions holds the complete logic of the artsparts site.
// The insitutions are organized over a slice.
type Institutions []*Institution

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
	ID          string                 `json:"id,omitempty" yaml:"id"`
	Name        string                 `json:"name,omitempty" yaml:"name"`
	Description string                 `json:"description,omitempty" yaml:"desc"`
	Collections map[string]*Collection `json:"collections,omitempty" yaml:"-"`
	Admins      []string               `json:"admins,omitempty" yaml:"admins"`
}

// NewInstitution takes a filepath and loads the configuration. Then
// it loads all the collections.
func NewInstitution(fpath string) (*Institution, error) {
	inst := &Institution{}
	confFile := filepath.Join(fpath, ConfFileName)
	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, inst)
	ls, err := ioutil.ReadDir(fpath)
	for _, d := range ls {
		if !d.IsDir() {
			continue
		}
		coll, err := NewCollection(filepath.Join(fpath, d.Name()))
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
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	License     string     `json:"license,omitempty"`
	Order       int        `json:"order,omitempty"`
	Artworks    []*Artwork `json:"artworks,omitempty"`
}

func NewCollection(fpath string) (*Collection, error) {
	return &Collection{}, nil
}

// Artwork is one element like for example a picture. The picture
// should be placed inside a folder. One artwork per folder! The
// foldername is then used as the id.
type Artwork struct {
	Timestamp   int    `json:"timestamp,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

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
