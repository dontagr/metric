package config

type Config struct {
	Log             Logging  `json:"Logging"`
	HTTPBindAddress string   `json:"address" env:"ADDRESS" validate:"required"`
	Security        Security `json:"Security"`
	PollInterval    int      `json:"poll_interval" env:"POLL_INTERVAL" validate:"required"`
	ReportInterval  int      `json:"report_interval" env:"REPORT_INTERVAL" validate:"required"`
	RateLimit       int      `json:"RateLimit" env:"RATE_LIMIT"`
	CryptoKey       string   `json:"crypto_key" env:"CRYPTO_KEY"`
}

type Security struct {
	Key string `json:"HashKey" env:"KEY"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}
