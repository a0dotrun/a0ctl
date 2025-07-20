package settings

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const EnvAccessToken = "A0_API_TOKEN"

func GetA0URL() string {
	settings, _ := ReadSettings()
	url := settings.GetBaseURL()
	if url == "" {
		url = settings.GetDefaultBaseURL()
	}
	return url
}

func GetA0HomeURL() string {
	settings, _ := ReadSettings()
	url := settings.GetHomeURL()
	if url == "" {
		url = settings.GetDefaultHomeURL()
	}
	return url
}

// TryToPersistChanges forces config changes to be written to disk.
func TryToPersistChanges() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to persist a0 settings file: %w", err)
	}
	return nil
}

// PersistChanges checks if there are any changes to the settings then persists them.
// More safer option
func PersistChanges() {
	if settings == nil || !settings.changed {
		return
	}

	if err := TryToPersistChanges(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
