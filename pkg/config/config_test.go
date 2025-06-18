package config

import (
	"os"
	"testing"
)

func TestConfig_IsTestFlag(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }() // Восстанавливаем оригинальное значение os.Args после тестов

	testCases := []struct {
		name         string
		args         []string
		expectedFlag bool
	}{
		{"No test flag", []string{"app"}, false},
		{"With test flag -test.v", []string{"app", "-test.v"}, true},
		{"With test flag -test.run", []string{"app", "-test.run=TestFunction"}, true},
		{"Non-test flag", []string{"app", "-prod"}, false},
		{"Multiple flags, including test", []string{"app", "-prod", "-test.v"}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args
			config := &Config{}
			got := config.IsTestFlag()
			if got != tc.expectedFlag {
				t.Errorf("For %s: expected %v, got %v", tc.name, tc.expectedFlag, got)
			}
		})
	}
}

type TestConfigData struct {
	HTTPBindAddress string `env:"HTTP_BIND_ADDRESS"`
	ReportInterval  int    `env:"REPORT_INTERVAL"`
	PollInterval    int    `env:"POLL_INTERVAL"`
}

func TestConfig_ReadFromEnv(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalHTTPBindAddress := os.Getenv("HTTP_BIND_ADDRESS")
	originalReportInterval := os.Getenv("REPORT_INTERVAL")
	originalPollInterval := os.Getenv("POLL_INTERVAL")
	defer func() {
		_ = os.Setenv("HTTP_BIND_ADDRESS", originalHTTPBindAddress)
		_ = os.Setenv("REPORT_INTERVAL", originalReportInterval)
		_ = os.Setenv("POLL_INTERVAL", originalPollInterval)
	}()

	testCases := []struct {
		name               string
		envHTTPBindAddress string
		envReportInterval  string
		envPollInterval    string
		expectedConfig     TestConfigData
		expectError        bool
	}{
		{
			name:               "All environment variables set",
			envHTTPBindAddress: "localhost:8080",
			envReportInterval:  "10",
			envPollInterval:    "5",
			expectedConfig:     TestConfigData{"localhost:8080", 10, 5},
			expectError:        false,
		},
		{
			name:               "Missing environment variable",
			envHTTPBindAddress: "",
			envReportInterval:  "10",
			envPollInterval:    "5",
			expectedConfig:     TestConfigData{"", 10, 5},
			expectError:        false,
		},
		{
			name:               "Invalid integer in environment variable",
			envHTTPBindAddress: "localhost:8080",
			envReportInterval:  "invalid",
			envPollInterval:    "5",
			expectedConfig:     TestConfigData{},
			expectError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Setenv("HTTP_BIND_ADDRESS", tc.envHTTPBindAddress)
			_ = os.Setenv("REPORT_INTERVAL", tc.envReportInterval)
			_ = os.Setenv("POLL_INTERVAL", tc.envPollInterval)

			configData := &TestConfigData{}
			cnf := &Config{Data: configData}
			err := cnf.ReadFromEnv()

			if tc.expectError {
				if err == nil {
					t.Errorf("%s: Expected an error but got nil", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("%s: Expected no error but got %v", tc.name, err)
				}
				if *configData != tc.expectedConfig {
					t.Errorf("%s: Expected config %+v, but got %+v", tc.name, tc.expectedConfig, *configData)
				}
			}
		})
	}
}
