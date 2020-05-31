package gopostal

import (
	"fmt"
	"github.com/ahmdrz/goinsta"
	"github.com/ahmdrz/goinsta/utilities"
	"github.com/yanatan16/golang-instagram/instagram"
	"os"
	"strconv"
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

func MakeEncodedCreds(username, password string) (string, error) {
	insta := goinsta.New(username, password)
	if err := insta.Login(); err != nil {
		return "", fmt.Errorf("Error logging in %v", err)
	}
	bytes, err := utilities.ExportAsBase64String(insta)
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

	privateAPI, err := utilities.ImportFromBase64String(options.EncodedCreds)
	if err != nil {
		return nil, fmt.Errorf("Error on importing creds")
	}

	//if err = privateAPI.Login(); err != nil {
	//	return nil, fmt.Errorf("Failed to login (%v)", err)
	//}

	// api := instagram.New(options.ClientID, options.ClientSecret, options.AccessToken, false)

	return &InstagramClient{
		Options:    options,
		privateAPI: privateAPI,
		// api:        api,
	}, nil
}

func (client *InstagramClient) Logout() error {
	return client.privateAPI.Logout()
}

func (client *InstagramClient) UploadPhoto(path, caption string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = client.privateAPI.UploadPhoto(file, caption,
		InstagramDefaultQuality, InstagramDefaultFilter)
	return err
}

func (client *InstagramClient) LastPhotoTime(userID string) (mostRecentTime time.Time, err error) {
	numericUserID, convErr := strconv.ParseInt(userID, 10, 64)
	if convErr != nil {
		err = convErr
		return
	}

	user := client.privateAPI.NewUser()
	user.ID = numericUserID

	feed := user.Feed()
	if !feed.Next() {
		err = fmt.Errorf("failed to sync feed")
		return
	}

	for _, media := range feed.Items {
		fmt.Printf("\nreading item %d\n", media.DeviceTimestamp)
		createdTime := time.Unix(0, media.DeviceTimestamp)
		if mostRecentTime.IsZero() || createdTime.After(mostRecentTime) {
			mostRecentTime = createdTime
		}
	}

	if mostRecentTime.IsZero() {
		err = fmt.Errorf("No media found.")
	}

	return
}
