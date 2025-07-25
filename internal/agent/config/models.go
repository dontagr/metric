package config

type Config struct {
	Log             Logging  `json:"Logging"`
	PollInterval    int      `json:"PollInterval" env:"POLL_INTERVAL" validate:"required"`
	ReportInterval  int      `json:"ReportInterval" env:"REPORT_INTERVAL" validate:"required"`
	HTTPBindAddress string   `json:"HTTPBindAddress" env:"ADDRESS" validate:"required"`
	Security        Security `json:"Security"`
	RateLimit       int      `json:"RateLimit" env:"RATE_LIMIT"`
}

type Security struct {
	Key string `json:"HashKey" env:"KEY"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}
