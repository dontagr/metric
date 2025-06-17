package config

type Config struct {
	PollInterval    int    `json:"PollInterval" env:"POLL_INTERVAL" validate:"required"`
	ReportInterval  int    `json:"ReportInterval" env:"REPORT_INTERVAL" validate:"required"`
	HTTPBindAddress string `json:"HTTPBindAddress" env:"ADDRESS" validate:"required"`
}
