package main

import (
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

func init() {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))

	store.MaxLength(math.MaxInt64)

	gothic.Store = store
	goth.UseProviders(
		twitter.New(
			conf.TwitterKey,
			conf.TwitterSecret,
			"http://localhost:3000/auth/twitter/callback"),
	)
}

func addAuthRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/twitter/callback", callbackHandler).Methods("GET")
	return r
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, r)
		return
	}
	fmt.Fprintf(w, "%#v", user)
}
