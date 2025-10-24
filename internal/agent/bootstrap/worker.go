package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/worker"
	"github.com/dontagr/metric/internal/agent/worker/pool"
)

var Worker = fx.Options(
	fx.Provide(
		worker.NewSender,
		worker.NewRefresher,
		pool.NewWPool,
	),
	fx.Invoke(
		func(sender *worker.Sender) {},
		func(sender *worker.Refresher) {},
	),
)
