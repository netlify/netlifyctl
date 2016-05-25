package auth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

var AccessToken string

func ClientCredentials() runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "netlify-go-cli")
		r.SetHeaderParam("Authorization", "Bearer "+chooseAccessToken())
		return nil
	})
}

func chooseAccessToken() string {
	if len(AccessToken) > 0 {
		return AccessToken
	}

	if token := loadAccessTokenFromFile(); len(token) > 0 {
		return token
	}

	return ""
}

func loadAccessTokenFromFile() string {
	home := os.Getenv("HOME")
	f, err := os.Open(filepath.Join(home, ".config", "netlify"))
	if err != nil {
		return ""
	}

	config := struct {
		AccessToken string
	}{}

	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return ""
	}

	return config.AccessToken
}
