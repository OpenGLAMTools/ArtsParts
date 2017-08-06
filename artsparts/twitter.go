package main

import (
	"image"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	artsparts "github.com/OpenGLAMTools/ArtsParts"
)

func initTwitter(conf Conf) {
	anaconda.SetConsumerKey(conf.Env["TWITTER_KEY"])
	anaconda.SetConsumerSecret(conf.Env["TWITTER_SECRET"])
}

type tweetResponse struct {
	twitterID       int64
	twitterIDString string
	mediaID         int64
	mediaIDString   string
}

func postPartTweet(ap *artsparts.Part, img image.Image, twitterAPI *anaconda.TwitterApi) error {
	log.Infoln("-----postPartTweet----")
	resp, err := tweetImage(ap.Text, img, twitterAPI)
	ap.TweetID = resp.twitterID
	ap.TweetIDString = resp.twitterIDString
	ap.MediaID = resp.mediaID
	ap.MediaIDString = resp.mediaIDString
	return err
}

func tweetImage(text string, img image.Image, twitterAPI *anaconda.TwitterApi) (tweetResponse, error) {
	var resp tweetResponse
	imgString, err := artsparts.ImageToBaseString(img)
	if err != nil {
		return resp, err
	}
	m, err := twitterAPI.UploadMedia(imgString)
	if err != nil {
		return resp, err
	}
	resp.mediaID = m.MediaID
	resp.mediaIDString = m.MediaIDString
	v := url.Values{
		"media_ids": []string{m.MediaIDString},
	}
	tweet, err := twitterAPI.PostTweet(text, v)
	log.Infof("TweetImage: %#v\n", tweet)
	resp.twitterID = tweet.Id
	resp.twitterIDString = tweet.IdStr
	return resp, err

}
