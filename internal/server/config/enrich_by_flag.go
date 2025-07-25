package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type FlagEnricher struct {
}

func (f *FlagEnricher) Process(cnf *Config) error {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	flagSet.Usage = cleanenv.FUsage(flagSet.Output(), cnf, nil, flagSet.Usage)

	serverAddrBind := flagSet.String("a", "", "bind addr http")
	storeInterval := flagSet.Int("i", 0, "time interval in seconds after which the current server readings are saved to disk")
	storeFilePath := flagSet.String("f", "", "path to the file where the current values are saved")
	storeRestore := flagSet.Bool("r", true, "load previously saved values for not for store")
	databaseDsn := flagSet.String("d", "", "string with the database connection address")
	key := flagSet.String("k", "", "key for encode with SHA256")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if *serverAddrBind != "" {
		cnf.HTTPServer.BindAddress = *serverAddrBind
	}

	if *databaseDsn != "" {
		cnf.DataBase.DatabaseDsn = *databaseDsn
		cnf.DataBase.Init = true
	} else if val, exists := os.LookupEnv(DatabaseDsn); exists && val != "" {
		cnf.DataBase.Init = true
	}

	_, exists := os.LookupEnv(EnvStoreInterval)
	if !exists && *storeInterval != 0 {
		cnf.Store.Interval = *storeInterval
	}

	_, exists = os.LookupEnv(EnvFileStoragePath)
	if !exists && *storeFilePath != "" {
		cnf.Store.FilePath = *storeFilePath
	}

	_, exists = os.LookupEnv(EnvRestore)
	if !exists {
		cnf.Store.Restore = *storeRestore
	}

	_, exists = os.LookupEnv(KEY)
	if !exists {
		cnf.Security.Key = *key
	}

	return nil
}
