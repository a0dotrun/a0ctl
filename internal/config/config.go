// Package config provides config management for CLI
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const AccessTokenEnv = "A0_API_TOKEN"

func GetA0URL() string {
	settings, _ := ReadSettings()
	url := settings.GetBaseURL()
	if url == "" {
		url = settings.GetDefaultBaseURL()
	}
	return url
}

func TryToPersistChanges() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to persist turso settings file: %w", err)
	}
	return nil
}

func PersistChanges() {
	if settings == nil || !settings.changed {
		return
	}

	if err := TryToPersistChanges(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
