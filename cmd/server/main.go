package main

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/bootstrap"
)

func main() {
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
