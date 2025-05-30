package config

var configFileNames = []string{"config.json"}

var configPaths = []string{"../../configs", "./configs"}

type Config struct {
	HTTPServer HTTPServer `json:"HttpServing"`
}

type HTTPServer struct {
	BindAddress string `json:"BindAddress" env:"ADDRESS" validate:"required"`
}
