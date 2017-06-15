package artsparts

// Institution defines a museum or another partner of the site,
// which offers collections.
//
// Following structure of the folders is needed:
//
// institution
// └───collections
//     └───collection1
//         ├───pic1
//         └───pic2
type Institution struct {
	Name        string
	Description string
	Collections map[string]*Collection
}

// Collection represents a group of artworks, which are presented
// together. It could be grouped after artist or a specific style.
type Collection struct {
	Name        string
	Description string
	License     string
	Artworks    []*Artwork
}

// Artwork is one element like for example a picture. The picture
// should be placed inside a folder. One artwork per folder! The
// foldername is then used as the id.
type Artwork struct {
	Timestamp   int
	Name        string
	Description string
}

type User struct {
	Twitter string
	Email   string
}
