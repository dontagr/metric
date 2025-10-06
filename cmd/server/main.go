package main

import (
	_ "net/http/pprof"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/bootstrap"
	start "github.com/dontagr/metric/pkg/service/print"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	start.PrintVersion(buildVersion, buildDate, buildCommit)
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
