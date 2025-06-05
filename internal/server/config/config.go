package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	if !isTestFlag() {
		enrichByFlag(&config)
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			fmt.Println(err)
		}
	}

	return &config, nil
}

func enrichByFlag(config *Config) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), config, nil, flagSet.Usage)

	serverAddrBind := flagSet.String("a", "", "bind addr http")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}
	if *serverAddrBind != "" {
		config.HTTPServer.BindAddress = *serverAddrBind
	}
}

func isTestFlag() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
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
