package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var serverAddrBind *string
var storeFilePath *string
var databaseDsn *string
var key *string
var cryptoKey *string
var storeInterval *int
var storeRestore *bool
var configShort *string
var config *string
var trustedSubnet *string

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
	storeInterval = flagSet.Int("i", 0, "time interval in seconds after which the current server readings are saved to disk")
	storeFilePath = flagSet.String("f", "", "path to the file where the current values are saved")
	storeRestore = flagSet.Bool("r", true, "load previously saved values for not for store")
	databaseDsn = flagSet.String("d", "", "string with the database connection address")
	key = flagSet.String("k", "", "key for encode with SHA256")
	cryptoKey = flagSet.String("crypto-key", "", "crypto-key for public")
	configShort = flagSet.String("c", "", "configuration file name")
	config = flagSet.String("config", "", "configuration file name")
	trustedSubnet = flagSet.String("t", "", "trusted subnet")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return nil
}

func (f *FlagEnricher) Process(cnf *Config) {
	if *serverAddrBind != "" {
		cnf.HTTPServer = *serverAddrBind
	}

	if *databaseDsn != "" {
		cnf.DatabaseDsn = *databaseDsn
		cnf.DataBase.Init = true
	} else if val, exists := os.LookupEnv(DatabaseDsn); exists && val != "" {
		cnf.DataBase.Init = true
	}

	_, exists := os.LookupEnv(EnvStoreInterval)
	if !exists && *storeInterval != 0 {
		cnf.Interval = *storeInterval
	}

	_, exists = os.LookupEnv(EnvFileStoragePath)
	if !exists && *storeFilePath != "" {
		cnf.FilePath = *storeFilePath
	}

	_, exists = os.LookupEnv(EnvRestore)
	if !exists {
		cnf.Restore = *storeRestore
	}

	_, exists = os.LookupEnv(KEY)
	if !exists {
		cnf.Security.Key = *key
	}

	_, exists = os.LookupEnv(CryptoKey)
	if !exists {
		cnf.CryptoKey = *cryptoKey
	}

	_, exists = os.LookupEnv(TrustedSubnet)
	if !exists {
		cnf.TrustedSubnet = *trustedSubnet
	}
}
