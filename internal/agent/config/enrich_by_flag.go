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
	reportInterval := flagSet.Int("r", 0, "report interval value in sec")
	pollInterval := flagSet.Int("p", 0, "poll interval value in sec")
	key := flagSet.String("k", "", "key for encode with SHA256")
	rateLimit := flagSet.Int("l", 0, "number of simultaneously outgoing requests to the server")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}
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

	return nil
}
