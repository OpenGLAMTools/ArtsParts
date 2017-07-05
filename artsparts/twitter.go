package main

import (
	"image"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

func initTwitter() {
	anaconda.SetConsumerKey(getenv("TWITTER_KEY"))
	anaconda.SetConsumerSecret(getenv("TWITTER_SECRET"))
}
func postTweetHandler(w http.ResponseWriter, r *http.Request) {
	/*	session, err := gothic.Store.Get(r, sessionName)
		if err != nil {
			log.Warningln("Error when calling session get():", err)
		}
		accesToken := session.Values["access_token"].(string)
		accesTokenSecret := session.Values["access_token_secret"].(string)
		twitterAPI := anaconda.NewTwitterApi(
			accesToken,
			accesTokenSecret,
		)
		err = postTweet("Text...", "../test/test.jpg", twitterAPI)
		if err != nil {
			log.Warningln("Error when tweeting: ", err)
		}*/
}
func postTweet(text string, img image.Image, twitterAPI *anaconda.TwitterApi) error {

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
