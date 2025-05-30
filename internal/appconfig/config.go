package appconfig

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
