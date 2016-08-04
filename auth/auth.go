package auth

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/viper"
)

func ClientCredentials() runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "netlifyctl")
		r.SetHeaderParam("Authorization", "Bearer "+chooseAccessToken())
		return nil
	})
}

func chooseAccessToken() string {
	// this will load it from the toml file, env, or cmdline
	accessToken := viper.GetString("access_token")
	return accessToken
}
