package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const EnvConfig = "CONFIG"

var serverAddrBind *string
var reportInterval *int
var pollInterval *int
var rateLimit *int
var key *string
var cryptoKey *string
var configShort *string
var config *string

type FlagEnricher struct {
}

func (f *FlagEnricher) GetFilePathAndName(paths []string, names []string) ([]string, []string) {
	if configShort == nil {
		return paths, names
	} else if *configShort != "" {
		names = append(names, *configShort)
	} else if *config != "" {
		names = append(names, *config)
	} else if name, exists := os.LookupEnv(EnvConfig); exists {
		names = append(names, name)
	}

	return paths, names
}

func (f *FlagEnricher) Init(cnf *Config) error {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), cnf, nil, flagSet.Usage)

	serverAddrBind = flagSet.String("a", "", "bind addr http")
	reportInterval = flagSet.Int("r", 0, "report interval value in sec")
	pollInterval = flagSet.Int("p", 0, "poll interval value in sec")
	key = flagSet.String("k", "", "key for encode with SHA256")
	rateLimit = flagSet.Int("l", 0, "number of simultaneously outgoing requests to the server")
	cryptoKey = flagSet.String("crypto-key", "", "crypto-key for public")
	configShort = flagSet.String("c", "", "configuration file name")
	config = flagSet.String("config", "", "configuration file name")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return nil
}

func (f *FlagEnricher) Process(cnf *Config) {
	if *reportInterval != 0 {
		cnf.ReportInterval = *reportInterval
	}
	if *pollInterval != 0 {
		cnf.PollInterval = *pollInterval
	}
	if *serverAddrBind != "" {
		cnf.HTTPBindAddress = *serverAddrBind
	}
	if *key != "" {
		cnf.Security.Key = *key
	}
	if *rateLimit != 0 {
		cnf.RateLimit = *rateLimit
	}
	if *cryptoKey != "" {
		cnf.CryptoKey = *cryptoKey
	}
}
