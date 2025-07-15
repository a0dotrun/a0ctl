// Package helpers provides utility functions for the a0ctl CLI tool.
package helpers

import (
	"fmt"
	"os"

	"github.com/a0dotrun/a0ctl/internal/config"

	"github.com/a0dotrun/a0ctl/internal/cli"
)

const accessTokenEnv = "A0_API_TOKEN"

var ErrNotLoggedIn = fmt.Errorf("user not logged in, please login with %s", cli.Emph("a0ctl auth login"))

func GetAccessToken() (string, error) {
	token, err := envAccessToken()
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}

	settings, err := config.ReadSettings()
	if err != nil {
		return "", fmt.Errorf("could not read token from settings file: %w", err)
	}

	token = settings.GetToken()
	if !IsJWTTokenValid(token) {
		return "", ErrNotLoggedIn
	}

	return token, nil
}

// envAccessToken retrieves the access token from the environment variable.
func envAccessToken() (string, error) {
	token := os.Getenv(accessTokenEnv)
	if token == "" {
		return "", nil
	}
	if !IsJWTTokenValid(token) {
		return "", fmt.Errorf("token in %s env var is invalid. Update the env var with a valid value, or unset it to use a token from the configuration file", accessTokenEnv)
	}
	return token, nil
}

func IsJWTTokenValid(token string) bool {
	if len(token) == 0 {
		return false
	}

	return false
}
