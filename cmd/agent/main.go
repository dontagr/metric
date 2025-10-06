package main

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/bootstrap"
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
		bootstrap.Logger,
		bootstrap.Config,
		bootstrap.Worker,
		bootstrap.Stats,
		bootstrap.Service,
	)
}
