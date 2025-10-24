package bootstrap

import (
	"go.uber.org/fx"

	agent "github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/pkg/config"
)

var Config = fx.Options(
	fx.Provide(newConfig),
)

func newConfig() (*agent.Config, error) {
	agentConfig := &agent.Config{}
	flagEnricher := &agent.FlagEnricher{}
	cnf := &config.Config{
		Data: agentConfig,
	}

	if !cnf.IsTestFlag() {
		err := flagEnricher.Init(agentConfig)
		if err != nil {
			return nil, err
		}
	}
	path, name := flagEnricher.GetFilePathAndName([]string{"../../../configs", "./configs"}, []string{"agent.json"})
	cnf.DefaultFilePaths = path
	cnf.DefaultFileNames = name

	cnf.ReadFromFile()
	if !cnf.IsTestFlag() {
		flagEnricher.Process(agentConfig)
	}

	err := cnf.ReadFromEnv()
	if err != nil {
		return nil, err
	}

	err = cnf.Validate()
	if err != nil {
		return nil, err
	}

	return agentConfig, err
}
