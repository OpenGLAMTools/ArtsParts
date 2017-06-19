package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	artsparts "github.com/OpenGLAMTools/ArtsParts"
	"github.com/disintegration/imaging"
	"github.com/markbates/goth/gothic"
)

func initTwitter() {
	anaconda.SetConsumerKey(getenv("TWITTER_KEY"))
	anaconda.SetConsumerSecret(getenv("TWITTER_SECRET"))
}
func postTweetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, sessionName)
	log.Println("Error when calling session get():", err)
	accesToken := session.Values["access_token"].(string)
	accesTokenSecret := session.Values["access_token_secret"].(string)
	twitterAPI := anaconda.NewTwitterApi(
		accesToken,
		accesTokenSecret,
	)
	err = postTweet("Text...", "../test/test.jpg", twitterAPI)
	if err != nil {
		log.Println("Error when tweeting: ", err)
	}
}
func postTweet(text, imagePath string, twitterAPI *anaconda.TwitterApi) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return err
	}
	imgString, err := artsparts.ImageToBaseString(img)
	if err != nil {
		return err
	}
	m, err := twitterAPI.UploadMedia(imgString)
	if err != nil {
		return err
	}
	v := url.Values{
		"media_ids": []string{m.MediaIDString},
	}
	_, err = twitterAPI.PostTweet(text, v)
	return err
}