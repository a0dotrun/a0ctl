package helpers

import "os"

func GetDaggerCloudToken() string {
	return os.Getenv("DAGGER_CLOUD_TOKEN")
}
