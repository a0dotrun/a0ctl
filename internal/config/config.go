// Package config provides config management for CLI
package config

const AccessTokenEnv = "A0_API_TOKEN"

func GetA0URL() string {
	settings, _ := ReadSettings()
	url := settings.GetBaseURL()
	if url == "" {
		url = settings.GetDefaultBaseURL()
	}
	return url
}
