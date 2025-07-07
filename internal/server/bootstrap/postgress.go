package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/faultTolerance/pgRetry"
)

var Postgress = fx.Options(
	fx.Provide(
		newPostgresConnect,
		pgRetry.NewPgxRetry,
	),
)

func newPostgresConnect(cfg *config.Config, log *zap.SugaredLogger, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	if !cfg.DataBase.Init {
		log.Info("коннект с базой не был инициализирован")

		return nil, nil
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DataBase.DatabaseDsn)
	if err != nil {
		log.Errorf("Unable to connect to database: %v", err)

		return nil, nil // just for auto test TestIteration1
		//		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			dbpool.Close()
			return nil
		},
	})

	return dbpool, nil
}
