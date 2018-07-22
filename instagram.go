package gopostal

import (
	"fmt"
	"github.com/ahmdrz/goinsta"
	"github.com/ahmdrz/goinsta/store"
	"github.com/yanatan16/golang-instagram/instagram"
	"net/url"
	"time"
)

const (
	InstagramDefaultQuality = 100
	InstagramDefaultFilter  = 0
)

type InstagramClientOptions struct {
	Username     string // Instagram username
	EncodedCreds string // Encoded instagram connection creds
	EncodeSecret string // Key for encoded creds
	ClientID     string // Api Client ID for reading
	ClientSecret string // Api Client Secret for reading
	AccessToken  string // Api Access Token
}

type InstagramClient struct {
	Options    *InstagramClientOptions
	privateAPI *goinsta.Instagram
	api        *instagram.Api
}

func MakeEncodedCreds(username, password, key string) (string, error) {
	insta := goinsta.New(username, password)
	if err := insta.Login(); err != nil {
		return "", fmt.Errorf("Error logging in %v", err)
	}
	bytes, err := store.Export(insta, []byte(key))
	if err != nil {
		return "", fmt.Errorf("Error on export")
	}

	insta.Logout()
	return string(bytes), nil
}

func NewInstagramClient(options *InstagramClientOptions) (*InstagramClient, error) {
	if options == nil {
		options = &InstagramClientOptions{}
	}

	if options.EncodeSecret == "" {
		return nil, fmt.Errorf("EncodeSecret must not be empty if you pass EncodedCreds")
	}
	privateAPI, err := store.Import([]byte(options.EncodedCreds), []byte(options.EncodeSecret))
	if err != nil {
		return nil, fmt.Errorf("Error on import")
	}

	if err = privateAPI.Login(); err != nil {
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
