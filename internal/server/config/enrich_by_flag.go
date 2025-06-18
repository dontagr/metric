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
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if *serverAddrBind != "" {
		cnf.HTTPServer.BindAddress = *serverAddrBind
	}

	return nil
}
