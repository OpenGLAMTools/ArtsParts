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
	"github.com/pkg/errors"
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
	s.HandleFunc("/logout", logoutHandler).Methods("GET")
	s.HandleFunc("", authHandler).Methods("GET")
	return r
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, sessionName)

	if err != nil {
		log.Warningln("logoutHandler: Error when session get():", err)
	}
	session.Values["userid"] = ""
	session.Values["twitter"] = ""
	session.Values["access_token"] = ""
	session.Values["access_token_secret"] = ""
	err = session.Save(r, w)
	if err != nil {
		log.Warningln("logoutHandler: Error on session save: ", err)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	session.Values["userid"] = gothUser.UserID
	session.Values["twitter"] = gothUser.NickName
	session.Values["access_token"] = gothUser.AccessToken
	session.Values["access_token_secret"] = gothUser.AccessTokenSecret
	err = session.Save(r, w)
	if err != nil {
		log.Warningln("Error when saving session: ", err)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getSessionValues(r *http.Request) map[string]string {
	session, err := gothic.Store.Get(r, sessionName)
	if err != nil {
		log.Error(errors.WithStack(err))
	}
	vals := make(map[string]string)
	vals["userid"] = getString("userid", session.Values)
	vals["twitter"] = getString("twitter", session.Values)
	vals["access_token"] = getString("access_token", session.Values)
	vals["access_token_secret"] = getString("access_token_secret", session.Values)
	return vals
}

func getString(key string, m map[interface{}]interface{}) string {
	i, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := i.(string)
	if !ok {
		return ""
	}
	return s
}
