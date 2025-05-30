package config

var configFileNames = []string{"config.json"}

var configPaths = []string{"../../configs", "./configs"}

type Config struct {
	HttpServing HttpServing `json:"HttpServing"`
}

type HttpServing struct {
	BindAddress string `json:"BindAddress" env:"ADDRESS" validate:"required"`
}
