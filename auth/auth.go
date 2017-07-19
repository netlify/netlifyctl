package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	AccessToken string

	validAuthPaths = [][]string{
		{".config", "netlify"},
		{".netlify", "config"},
	}
)

func ClientCredentials() runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		token := chooseAccessToken()
		if token == "" {
			return errors.New("No access token found. Please login.")
		}

		r.SetHeaderParam("User-Agent", "netlifyctl")
		r.SetHeaderParam("Authorization", "Bearer "+token)
		return nil
	})
}

func NoCredentials() runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "netlifyctl")
		return nil
	})
}

func SaveToken(token string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	args := append([]string{home}, validAuthPaths[0]...)
	f, err := os.OpenFile(filepath.Join(args...), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	config := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	return json.NewEncoder(f).Encode(&config)
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
	home, err := homedir.Dir()
	if err != nil {
		return ""
	}

	var f *os.File
	for _, p := range validAuthPaths {
		args := append([]string{home}, p...)
		f, err = os.Open(filepath.Join(args...))
		if err == nil {
			break
		}
	}

	if err != nil || f == nil {
		return ""
	}
	defer f.Close()

	config := struct {
		AccessToken string `json:"access_token"`
	}{}

	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return ""
	}

	return config.AccessToken
}
