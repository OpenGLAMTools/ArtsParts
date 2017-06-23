package main

import (
	"math"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

var sessionName = "ap-user-session"
var store *sessions.FilesystemStore

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
		log.Errorf("No Env value set for %s", key)
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
		session, err := gothic.Store.Get(r, sessionName)
		if err != nil {
			log.Warningln("Error when session get():", err)
		}
		//session.Values["gothUser"] = gothUser
		session.Values["userid"] = gothUser.UserID
		session.Values["twitter"] = gothUser.NickName
		session.Values["access_token"] = gothUser.AccessToken
		session.Values["access_token_secret"] = gothUser.AccessTokenSecret
		err = session.Save(r, w)
		if err != nil {
			log.Warningln("Error on session save: ", err)
		}
	} else {
		gothic.BeginAuthHandler(w, r)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Warningln("Error completing auth: ", err)
		return
	}
	session, err := gothic.Store.Get(r, sessionName)
	if err != nil {
		log.Warningln("Error when calling session get():", err)
	}
	session.Values["gothUser"] = gothUser
	err = session.Save(r, w)
	if err != nil {
		log.Warningln("Error when saving session: ", err)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
