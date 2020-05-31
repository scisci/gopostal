package gopostal

import (
	"github.com/ahmdrz/goinsta/utilities"
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

	var encodedString string

	{
		creds, err := MakeEncodedCreds(username, password)

		if err != nil {
			t.Errorf("Failed to make encoded creds %v", err)
			return
		}
		encodedString = creds
	}

	time.Sleep(3 * time.Second)

	{
		privateAPI, err := utilities.ImportFromBase64String(encodedString)

		if err != nil {
			t.Errorf("Failed to import from encoded %v", err)
		}

		privateAPI.Logout()
	}

	log.Printf("Success %s\n", encodedString)
}
