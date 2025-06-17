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
	agentConfig := &server.Config{}
	flagEnricher := &server.FlagEnricher{}
	cnf := &config.Config{
		Data:             agentConfig,
		DefaultFilePaths: []string{"../../../configs", "./configs"},
		DefaultFileNames: []string{"server.json"},
	}

	cnf.ReadFromFile()
	if !cnf.IsTestFlag() {
		err := flagEnricher.Process(agentConfig)
		if err != nil {
			return nil, err
		}
	}

	err := cnf.ReadFromEnv()
	if err != nil {
		return nil, err
	}

	err = cnf.Validate()
	if err != nil {
		return nil, err
	}

	return agentConfig, nil
}
