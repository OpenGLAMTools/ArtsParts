package artsparts

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/OpenGLAMTools/ArtsParts/helpers"
	"github.com/OpenGLAMTools/ArtsParts/shortlink"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

// Artwork is one element like for example a picture. The picture
// should be placed inside a folder. One artwork per folder! The
// foldername is then used as the id.
// The conf file here is stored as JSON, because the content is created
// and edited via a configuration dialog.
// Path contains the part of the url how the artwork can be found:
// /[institution]/[collection]/[artwork]
type Artwork struct {
	Timestamp     string   `json:"timestamp"`
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	TweetID       int64    `json:"tweet,omitempty"`
	TweetIDString string   `json:"tweet_str,omitempty"`
	HashTag       string   `json:"hashtag,omitempty"`
	URIPath       string   `json:"-"`
	ShortLink     string   `json:"shortlink,omitempty"`
	Parts         []*Part  `json:"Parts,omitempty"`
	Meta          MetaData `json:"meta"`
	fpath         string
	collection    *Collection
}

type MetaData struct {
	Artist string `json:"artist"`
	Date   string `json:"date"`
	Link   string `json:"link"`
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
			// file exists but another error occurs
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
	artw.ShortLink, err = shortlink.GetShort("/artwork" + artw.URIPath)
	return artw, err
}

// License returns the license of the collection. Method is used inside the templates
func (a *Artwork) License() string {
	return a.collection.License
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

// AddPart adds a Part to the artwork
func (artw *Artwork) AddPart(p *Part) {
	artw.Parts = append(artw.Parts, p)
}

// Artpart creates the part of the artwork. Every number is relative to the size
// of the picture. To get the x value in pixel you need to multiply it to the width
func (artw *Artwork) Artpart(p *Part) (image.Image, error) {
	img, err := artw.Image()
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	imgWidth := float64(bounds.Max.X)
	imgHeigth := float64(bounds.Max.Y)

	rect := image.Rect(
		int(p.X*imgWidth),
		int(p.Y*imgHeigth),
		int((p.X+p.Width)*imgWidth),
		int((p.Y+p.Height)*imgHeigth))
	artp := imaging.Crop(img, rect)
	return artp, nil
}

// Image returns the Image of the artwork
func (artw *Artwork) Image() (image.Image, error) {
	f, err := artw.ImgFile()
	if err != nil {
		return nil, err
	}
	return imaging.Open(filepath.Join(artw.fpath, f))
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

// InstitutionTwitter returns the twitter name of the institution
func (artw *Artwork) InstitutionTwitter() string {
	return artw.collection.institution.TwitterName
}

// IsAdminUser returns true, when the given user has admin rights
func (artw *Artwork) IsAdminUser(userName string) bool {
	return helpers.StringInSlice(userName, artw.collection.institution.Admins)
}

func (artw *Artwork) dataFilePath() string {
	return filepath.Join(artw.fpath, DataFileName)
}
