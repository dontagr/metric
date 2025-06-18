package config

type Config struct {
	HTTPServer HTTPServer `json:"HttpServing"`
}

type HTTPServer struct {
	BindAddress string `json:"BindAddress" env:"ADDRESS" flag:"a" validate:"required"`
}
