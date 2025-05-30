package config

var configFileNames = []string{"agent.json"}

var configPaths = []string{"../../configs", "./configs"}

type Config struct {
	PollInterval   int `json:"PollInterval"`
	ReportInterval int `json:"ReportInterval"`
}
