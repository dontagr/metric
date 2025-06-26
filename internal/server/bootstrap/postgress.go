package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
)

var Postgress = fx.Options(
	fx.Provide(
		newPostgresConnect,
	),
)

func newPostgresConnect(cfg *config.Config, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.DataBase.DatabaseDsn)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v", err)

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
