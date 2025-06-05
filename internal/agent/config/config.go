package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig() (*Config, error) {
	return loadConfig()
}

func loadConfig() (*Config, error) {
	var config Config
	for _, file := range getConfigFile() {
		absPath, _ := filepath.Abs(file)
		fmt.Printf("Read config from path: %s\n", absPath)
		err := cleanenv.ReadConfig(absPath, &config)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Reading was successful")
		}
	}

	validate := validator.New()
	validateError := validate.Struct(&config)

	if validateError != nil {
		return nil, validateError
	}

	enrichByFlag(&config)

	return &config, nil
}

func enrichByFlag(config *Config) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	serverAddrBind := flagSet.String("a", "", "bind addr http")
	reportInterval := flagSet.Int("r", 0, "report interval value in sec")
	pollInterval := flagSet.Int("p", 0, "poll interval value in sec")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}
	if *reportInterval != 0 {
		config.ReportInterval = *reportInterval
	}
	if *pollInterval != 0 {
		config.PollInterval = *pollInterval
	}
	if *serverAddrBind != "" {
		config.HTTPBindAddress = *serverAddrBind
	}
}

func getConfigFile() []string {
	files := make([]string, 0)
	for _, path := range configPaths {
		for _, fileName := range configFileNames {
			files = append(files, fmt.Sprintf("%s/%s", path, fileName))
		}
	}

	return files
}
