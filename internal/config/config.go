// Package config provides config management for CLI
package config

func GetA0URL() string {
	settings, _ := ReadSettings()
	url := settings.GetBaseURL()
	if url == "" {
		url = settings.GetDefaultBaseURL()
	}
	return url
}
