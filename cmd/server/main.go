package main

import (
	"fmt"
	_ "net/http/pprof"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/bootstrap"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: \"%s\"\n", buildVersion)
	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build date: \"%s\"\n", buildDate)
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build commit: \"%s\"\n", buildCommit)

	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		bootstrap.Postgres,
		bootstrap.Store,
		bootstrap.Config,
		bootstrap.Server,
		bootstrap.Route,
		bootstrap.Service,
		bootstrap.Logger,
		bootstrap.Filer,
	)
}
