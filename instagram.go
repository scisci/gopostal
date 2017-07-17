package gopostal

import (
	"fmt"
	"github.com/ahmdrz/goinsta"
	"github.com/yanatan16/golang-instagram/instagram"
	"net/url"
	"time"
)

const (
	InstagramDefaultQuality = 100
	InstagramDefaultFilter  = 0
)

type InstagramClientOptions struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

type InstagramClient struct {
	Options    *InstagramClientOptions
	privateAPI *goinsta.Instagram
	api        *instagram.Api
}

func NewInstagramClient(options *InstagramClientOptions) (*InstagramClient, error) {
	if options == nil {
		options = &InstagramClientOptions{}
	}

	privateAPI := goinsta.New(options.Username, options.Password)
	if err := privateAPI.Login(); err != nil {
		return nil, err
	}

	api := instagram.New(options.ClientID, options.ClientSecret, options.AccessToken, false)

	return &InstagramClient{
		Options:    options,
		privateAPI: privateAPI,
		api:        api,
	}, nil
}

func (client *InstagramClient) Logout() error {
	return client.privateAPI.Logout()
}

func (client *InstagramClient) UploadPhoto(path, caption string) error {
	return client.UploadPhotoWithUploadID(path, caption, client.privateAPI.NewUploadID())
}

func (client *InstagramClient) UploadPhotoWithUploadID(path, caption string, uploadID int64) error {
	_, err := client.privateAPI.UploadPhoto(path, caption, uploadID,
		InstagramDefaultQuality, InstagramDefaultFilter)
	return err
}

func (client *InstagramClient) LastPhotoTime(userID string) (mostRecentTime time.Time, err error) {
	params := url.Values{}
	params.Set("count", "1")
	resp, err := client.api.GetUserRecentMedia(userID, params)
	if err != nil {
		return
	}

	for _, media := range resp.Medias {
		if createdTime, timeErr := media.CreatedTime.Time(); timeErr != nil {
			continue
		} else if mostRecentTime.IsZero() || createdTime.After(mostRecentTime) {
			mostRecentTime = createdTime
		}
	}

	if mostRecentTime.IsZero() {
		err = fmt.Errorf("No media found.")
	}

	return
}
