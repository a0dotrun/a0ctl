package api

import (
	"fmt"
	"os"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/settings"
)

// IsJWTTokenValid validates token.
func IsJWTTokenValid(token string) bool {
	if len(token) == 0 {
		return false
	}

	client, err := MakeClient(token)
	if err != nil {
		return false
	}

	r, err := client.Tokens.Validate()
	if err != nil {
		return false
	}

	return r
}

var ErrNotLoggedIn = fmt.Errorf(
	"user not logged in, please login with %s", cli.Emph("a0ctl auth login"))

func GetAccessToken() (string, error) {
	token, err := envAccessToken()
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}

	// env has no token, read from config.
	// env variable takes precedence over config file.
	config, err := settings.ReadSettings()
	if err != nil {
		return "", fmt.Errorf("could not read token from settings file: %w", err)
	}

	token = config.GetToken()
	if !IsJWTTokenValid(token) {
		return "", ErrNotLoggedIn
	}

	return token, nil
}

// envAccessToken retrieves the access token from the environment variable.
func envAccessToken() (string, error) {
	token := os.Getenv(settings.EnvAccessToken)
	if token == "" {
		return "", nil
	}
	if !IsJWTTokenValid(token) {
		return "", fmt.Errorf("token in %s env var is invalid. Update the env var with a valid value, or unset it to use a token from the configuration file", settings.EnvAccessToken)
	}
	return token, nil
}
