package artsparts

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
	ID          string
	Name        string
	Description string
	Collections map[string]*Collection
	Admins      []string
}

// Collection represents a group of artworks, which are presented
// together. It could be grouped after artist or a specific style.
// The Order property allows to sort the collections of a institution.
type Collection struct {
	ID          string
	Name        string
	Description string
	License     string
	Order       int
	Artworks    []*Artwork
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
