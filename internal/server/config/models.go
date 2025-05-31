package config

var configFileNames = []string{"server.json"}

var configPaths = []string{"../../../configs", "./configs"}

type Config struct {
	HTTPServer HTTPServer `json:"HttpServing"`
}

type HTTPServer struct {
	BindAddress string `json:"BindAddress" env:"ADDRESS" validate:"required"`
}
