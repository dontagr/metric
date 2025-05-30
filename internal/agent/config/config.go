package config

import (
	"fmt"
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

	return &config, nil
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
