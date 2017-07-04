// Package shortlink is a very simple link shortener.
package shortlink

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// DbFile is the filename of the bolt db file
var DbFile = "links.db"

// URIPrefix is the prefix for the router
var URIPrefix = "/s"

// Link represents a link with an id
type Link struct {
	ID  int    `storm:"id,increment"`
	URL string `storm:"unique"`
}

// GetShort returns the short url for a url
func GetShort(url string) (string, error) {
	id, err := GetID(url)
	return fmt.Sprintf("%s/%d", URIPrefix, id), err
}

// GetID returns a unique id to a given url. If a url don't have a ID a new
// is created.
func GetID(url string) (int, error) {
	db, err := storm.Open(DbFile)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer db.Close()

	l := Link{URL: url}
	err = db.One("URL", url, &l)
	if err == storm.ErrNotFound {
		err = db.Save(&l)
	}
	return l.ID, errors.WithStack(err)
}

func AddRoute(r *mux.Router) *mux.Router {
	r.HandleFunc(URIPrefix+"/{id}", Redirect)
	return r
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	db, err := storm.Open(DbFile)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	defer db.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	var l Link
	err = db.One("ID", id, &l)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File not found"))
		return
	}
	http.Redirect(w, r, l.URL, http.StatusTemporaryRedirect)
}
