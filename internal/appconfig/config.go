package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AppName string `json:"name"`
	Region  string `json:"region"`
}

func NewConfig(name, region string) Config {
	return Config{
		AppName: name,
		Region:  region,
	}
}

func ResolveAppConfig(cwd string) (Config, error) {
	configDir := filepath.Join(cwd, ".a0")
	configFile := filepath.Join(configDir, "app.json")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return Config{}, fmt.Errorf(".a0 directory not found in current directory")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("app.json not found in .a0 directory")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read app.json: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse app.json: %w", err)
	}

	if config.AppName == "" {
		return Config{}, fmt.Errorf("name field is required in app.json")
	}

	if config.Region == "" {
		return Config{}, fmt.Errorf("region field is required in app.json")
	}

	return config, nil
}
