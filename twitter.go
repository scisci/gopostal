package gopostal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"
)

type TwitterClientOptions struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type TwitterClient struct {
	Options *TwitterClientOptions
	api     *anaconda.TwitterApi
}

var noMediaFoundErr = errors.New("No media found.")

func NewTwitterClient(options *TwitterClientOptions) (*TwitterClient, error) {
	if options == nil {
		options = &TwitterClientOptions{}
	}

	anaconda.SetConsumerKey(options.ConsumerKey)
	anaconda.SetConsumerSecret(options.ConsumerSecret)
	api := anaconda.NewTwitterApi(options.AccessToken, options.AccessSecret)

	return &TwitterClient{
		Options: options,
		api:     api,
	}, nil
}

func (client *TwitterClient) IsNoMediaFoundErr(err error) bool {
	return err == noMediaFoundErr
}

func (client *TwitterClient) UploadPhoto(path, caption string) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	mediaResponse, err := client.api.UploadMedia(base64.StdEncoding.EncodeToString(data))
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("media_ids", strconv.FormatInt(mediaResponse.MediaID, 10))

	_, err = client.api.PostTweet(caption, v)
	if err != nil {
		return err
	}

	return nil
}

func (client *TwitterClient) LastPhotoTime(userID string) (mostRecentTime time.Time, err error) {
	params := url.Values{}
	params.Set("count", "10")
	params.Set("user_id", userID)

	tweets, err := client.api.GetUserTimeline(params)
	if err != nil {
		return
	}

	for _, tweet := range tweets {
		if len(tweet.Entities.Media) == 0 {
			continue
		}

		if createdTime, timeErr := tweet.CreatedAtTime(); timeErr != nil {
			fmt.Printf("error getting created at time for tweet: (%v)\n", timeErr)
			continue
		} else if mostRecentTime.IsZero() || createdTime.After(mostRecentTime) {
			mostRecentTime = createdTime
		}
	}

	if mostRecentTime.IsZero() {
		err = noMediaFoundErr
	}

	return
}
