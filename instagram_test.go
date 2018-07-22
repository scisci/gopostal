package gopostal

import (
	"github.com/ahmdrz/goinsta/store"
	"log"
	"os"
	"testing"
	"time"
)

func TestExportImport(t *testing.T) {
	username := os.Getenv("INSTA_USERNAME")
	password := os.Getenv("INSTA_PASSWORD")
	if len(username)*len(password) == 0 && os.Getenv("INSTA_PULL") != "true" {
		t.Skip("Username or Password is empty")
	}

	var key = "RH1tCpR80AQ3WzXJ" //32byte key for AES
	var encodedString string

	{
		creds, err := MakeEncodedCreds(username, password, key)

		if err != nil {
			t.Error(err)
			return
		}
		encodedString = creds
	}

	time.Sleep(3 * time.Second)

	{
		privateAPI, err := store.Import([]byte(encodedString), []byte(key))

		if err != nil {
			t.Error(err)
		}

		privateAPI.Logout()
	}

	log.Printf("Success %s\n", encodedString)
}
