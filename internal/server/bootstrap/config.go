package bootstrap

import (
	"go.uber.org/fx"

	server "github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/pkg/config"
)

var Config = fx.Options(
	fx.Provide(newConfig),
)

func newConfig() (*server.Config, error) {
	serverConfig := &server.Config{}
	flagEnricher := &server.FlagEnricher{}
	cnf := &config.Config{
		Data:             serverConfig,
		DefaultFilePaths: []string{"../../../configs", "./configs"},
		DefaultFileNames: []string{"server.json"},
	}

	if !cnf.IsTestFlag() {
		err := flagEnricher.Init(serverConfig)
		if err != nil {
			return nil, err
		}
	}
	path, name := flagEnricher.GetFilePathAndName([]string{"../../../configs", "./configs"}, []string{"server.json"})
	cnf.DefaultFilePaths = path
	cnf.DefaultFileNames = name

	cnf.ReadFromFile()
	if !cnf.IsTestFlag() {
		flagEnricher.Process(serverConfig)
	}

	err := cnf.ReadFromEnv()
	if err != nil {
		return nil, err
	}

	err = cnf.Validate()
	if err != nil {
		return nil, err
	}

	return serverConfig, nil
}
