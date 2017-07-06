package main

import (
	"image"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

func initTwitter() {
	anaconda.SetConsumerKey(getenv("TWITTER_KEY"))
	anaconda.SetConsumerSecret(getenv("TWITTER_SECRET"))
}

func postTweet(ap *artsparts.Part, img image.Image, twitterAPI *anaconda.TwitterApi) error {

	imgString, err := artsparts.ImageToBaseString(img)
	if err != nil {
		return err
	}
	m, err := twitterAPI.UploadMedia(imgString)
	if err != nil {
		return err
	}
	ap.MediaID = m.MediaID
	v := url.Values{
		"media_ids": []string{m.MediaIDString},
	}
	tweet, err := twitterAPI.PostTweet(ap.Text, v)
	if err != nil {
		return err
	}
	ap.TweetID = tweet.Id
	return nil
}
