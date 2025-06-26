package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
)

var Postgress = fx.Options(
	fx.Provide(
		newPostgresConnect,
	),
)

func newPostgresConnect(cfg *config.Config, lc fx.Lifecycle) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DataBase.DatabaseDsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close(ctx)
		},
	})

	return conn, nil
}
