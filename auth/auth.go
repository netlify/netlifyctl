package auth

import "github.com/spf13/viper"

func AuthToken() string {
	accessToken := viper.GetString("access_token")
	return accessToken
}
