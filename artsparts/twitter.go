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

func postPartTweet(ap *artsparts.Part, img image.Image, twitterAPI *anaconda.TwitterApi) error {
	log.Infoln("-----postPartTweet----")
	twitterID, mediaID, err := tweetImage(ap.Text, img, twitterAPI)
	ap.TweetID = twitterID
	ap.MediaID = mediaID
	return err
}

func tweetImage(text string, img image.Image, twitterAPI *anaconda.TwitterApi) (twitterID, mediaID int64, err error) {
	imgString, err := artsparts.ImageToBaseString(img)
	if err != nil {
		return twitterID, mediaID, err
	}
	m, err := twitterAPI.UploadMedia(imgString)
	if err != nil {
		return twitterID, mediaID, err
	}
	mediaID = m.MediaID
	v := url.Values{
		"media_ids": []string{m.MediaIDString},
	}
	tweet, err := twitterAPI.PostTweet(text, v)
	log.Infof("TweetImage: %#v\n", tweet)
	twitterID = tweet.Id
	return twitterID, mediaID, err

}
