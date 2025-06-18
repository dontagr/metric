package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/worker"
)

var Worker = fx.Options(
	fx.Provide(
		worker.NewSender,
		worker.NewRefresher,
	),
	fx.Invoke(
		func(sender *worker.Sender) {},
		func(sender *worker.Refresher) {},
	),
)
