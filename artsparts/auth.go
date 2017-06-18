package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

func initAuth(conf Conf) {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte(conf.SessionSecret))
	store.MaxLength(math.MaxInt64)
	gothic.Store = store
	goth.UseProviders(
		twitter.New(
			getenv("TWITTER_KEY"),
			getenv("TWITTER_SECRET"),
			"http://localhost:3000/auth/twitter/callback"),
	)
}

func getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("No Env value set for %s", key)
	}
	return val
}

func addAuthRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/auth/{provider}").Subrouter()
	s.HandleFunc("/callback", authCallbackHandler).Methods("GET")
	s.HandleFunc("", authHandler).Methods("GET")
	return r
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		fmt.Fprintf(w, "UserFound\n\n%#v", gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, r)
		return
	}
	fmt.Fprintf(w, "%#v", user)
}
