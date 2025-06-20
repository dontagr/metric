package config

type Config struct {
	Log        Logging    `json:"Logging"`
	HTTPServer HTTPServer `json:"HttpServing"`
}

type HTTPServer struct {
	BindAddress string `json:"BindAddress" env:"ADDRESS" flag:"a" validate:"required"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}
