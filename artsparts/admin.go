package main

import (
	"net/http"

	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/OpenGLAMTools/ArtsParts/helpers"
)

// Admin is used for managing all admin issues. It also implements
// the http.Handler interface and is used to serve all the admin
// pages.
type Admin struct {
	institutions artsparts.Institutions
}

// NewAdmin returns a pointer to a new instance
func NewAdmin(i artsparts.Institutions) *Admin {
	return &Admin{
		institutions: i,
	}
}

// Institutions takes the twitterName of the user and returns the institutions
// where he/she is admin of
func (a *Admin) Institutions(twitterName string) artsparts.Institutions {
	out := artsparts.Institutions{}
	for _, i := range a.institutions {
		if helpers.StringInSlice(twitterName, i.Admins) {
			out = append(out, i)
		}
	}
	return out
}

func (a *Admin) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
