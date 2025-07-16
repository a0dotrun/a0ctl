package api

import (
	"fmt"
	"os"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/config"
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

	// env has no token, read from settings.
	// env variable takes precedence over settings file.
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
	token := os.Getenv(config.EnvAccessToken)
	if token == "" {
		return "", nil
	}
	if !IsJWTTokenValid(token) {
		return "", fmt.Errorf("token in %s env var is invalid. Update the env var with a valid value, or unset it to use a token from the configuration file", config.EnvAccessToken)
	}
	return token, nil
}
