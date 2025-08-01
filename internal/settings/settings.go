// Package settings provides config management for CLI
package settings

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sync"

	"github.com/a0dotrun/a0ctl/internal/cli"
	"github.com/a0dotrun/a0ctl/internal/flags"
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

const (
	a0DefaultBaseURL = "https://api.a0.run"
	a0DefaultHomeURL = "https://a0.run"
)

type Settings struct {
	changed bool
}

var (
	settings *Settings
	mu       sync.Mutex
)

func ReadSettings() (*Settings, error) {
	mu.Lock()
	defer mu.Unlock()

	if settings != nil {
		return settings, nil
	}

	settings = &Settings{}

	configPath := configdir.LocalConfig("a0")
	err := viper.BindEnv("config-path", "A0_CONFIG_PATH")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("baseURL", "A0_API_BASEURL")
	if err != nil {
		return nil, err
	}

	err = viper.BindEnv("homeURL", "A0_HOME_BASEURL")
	if err != nil {
		return nil, err
	}

	configPathFlag := viper.GetString("config-path")
	if len(configPathFlag) > 0 {
		configPath = configPathFlag
	}

	err = configdir.MakePath(configPath)
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	configFile := path.Join(configPath, "settings.json")
	if abs, err := filepath.Abs(configFile); err == nil {
		configFile = abs
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		var configParseError viper.ConfigParseError
		switch {
		case errors.As(err, &configFileNotFoundError):
			// Force config creation
			if err := viper.SafeWriteConfig(); err != nil {
				return nil, err
			}
		case errors.As(err, &configParseError):
			if flags.ResetConfig() {
				err := viper.WriteConfig()
				if err != nil {
					return nil, err
				}
				break
			}
			warning := cli.Warn("Warning")
			// FIXME: requires implementation
			flag := cli.Emph("--reset-config")
			fmt.Printf("%s: could not parse JSON config from file %s\n", warning, cli.Emph(configFile))
			fmt.Printf("Fix the syntax errors on the file, or use the %s flag to replace it with a fresh one.\n", flag)
			fmt.Printf("E.g. a0ctl auth login --reset-config\n")
			return nil, err
		default:
			return nil, err
		}
	}

	return settings, nil
}

func (s *Settings) GetToken() string {
	return viper.GetString("token")
}

func (s *Settings) GetBaseURL() string {
	return viper.GetString("baseURL")
}

func (s *Settings) GetDefaultBaseURL() string {
	return a0DefaultBaseURL
}

func (s *Settings) GetHomeURL() string {
	return viper.GetString("homeURL")
}

func (s *Settings) GetDefaultHomeURL() string {
	return a0DefaultHomeURL
}

func (s *Settings) GetUsername() string {
	return viper.GetString("username")
}

func (s *Settings) SetToken(token string) {
	viper.Set("token", token)
	s.changed = true
}

func (s *Settings) SetUsername(username string) {
	viper.Set("username", username)
	s.changed = true
}
