package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/dontagr/metric/models"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS "metric" (
	  "id" VARCHAR(255) not null,
	  "mtype" varchar(255) not null,
	  "delta" BIGINT null,
	  "value" DOUBLE PRECISION null,
	  "hash" VARCHAR(255) not null,
	  constraint "metric_pkey" primary key ("id", "mtype")
	)`
	insertSql = `INSERT INTO metric (id, mtype, delta, value, hash) 
	  VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, mtype) DO UPDATE SET
	  mtype = EXCLUDED.mtype,
	  delta = EXCLUDED.delta,
	  value = EXCLUDED.value,
	  hash = EXCLUDED.hash`
	searchSql = `SELECT id, mtype, delta, value, hash 
      FROM metric 
      WHERE id=$1 AND mtype=$2`
	selectAllSql = `SELECT id, mtype, delta, value, hash FROM metric`
	truncateSql  = `TRUNCATE TABLE metric`
)

type pg struct {
	dbpool *pgxpool.Pool
	name   string
}

func RegisterStorePG(ms *StoreFactory, dbpool *pgxpool.Pool, lc fx.Lifecycle) {
	pg := pg{
		dbpool: dbpool,
		name:   models.StorePg,
	}

	if dbpool != nil {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return pg.addShema(ctx)
			},
		})
	}

	ms.SetStory(&pg)
}

func (pg *pg) GetName() string {
	return pg.name
}

func (pg *pg) addShema(ctx context.Context) error {
	_, err := pg.dbpool.Exec(ctx, createTable)
	return err
}

func (pg *pg) LoadMetric(id string, mType string) *models.Metrics {
	var metrics models.Metrics

	err := pg.dbpool.QueryRow(
		context.Background(),
		searchSql,
		id, mType,
	).Scan(&metrics.ID, &metrics.MType, &metrics.Delta, &metrics.Value, &metrics.Hash)
	if err != nil {
		fmt.Printf("Загрузка не удалась для (id: %s, mtype: %s): %v\n", id, mType, err)

		return &models.Metrics{}
	}

	return &metrics
}

func (pg *pg) SaveMetric(metrics *models.Metrics) {
	id, mtype, delta, value, hash := pg.unpack(metrics)

	_, err := pg.dbpool.Exec(context.Background(), insertSql, id, mtype, delta, value, hash)
	if err != nil {
		fmt.Printf("Ошибка при сохранении метрики: %v\n", err)
	}
}

func (pg *pg) unpack(metrics *models.Metrics) (string, string, *int64, *float64, string) {
	return metrics.ID, metrics.MType, metrics.Delta, metrics.Value, metrics.Hash
}

func (pg *pg) ListMetric() map[string]*models.Metrics {
	r := make(map[string]*models.Metrics)

	rows, err := pg.dbpool.Query(context.Background(), selectAllSql)
	if err != nil {
		fmt.Printf("Ошибка при извлечении метрик: %v\n", err)
		return r
	}
	defer rows.Close()

	for rows.Next() {
		var metrics models.Metrics
		err := rows.Scan(&metrics.ID, &metrics.MType, &metrics.Delta, &metrics.Value, &metrics.Hash)
		if err != nil {
			fmt.Printf("Ошибка при сканировании метрики: %v\n", err)
			continue
		}
		r[fmt.Sprintf("%s_%s", metrics.MType, metrics.ID)] = &metrics
	}

	return r
}

func (pg *pg) RestoreMetricCollection(collection map[string]*models.Metrics) {
	tx, txErr := pg.dbpool.Begin(context.Background())
	if txErr != nil {
		fmt.Printf("Ошибка начала транзакции: %v\n", txErr)
		return
	}
	defer func(txErr *error, count int) {
		if *txErr != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				fmt.Printf("Ошибка отката транзакции: %v\n", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(context.Background()); commitErr != nil {
				fmt.Printf("Ошибка при коммите транзакции: %v\n", commitErr)
			}
			fmt.Printf("\u001B[032mДанные в базе востановлены, всего метрик: %d\u001B[0m\n", count)
		}
	}(&txErr, len(collection))

	_, execErr := tx.Exec(context.Background(), truncateSql)
	if execErr != nil {
		fmt.Printf("Ошибка при восстановлении метрики: %v\n", execErr)
		txErr = execErr
		return
	}

	for _, metrics := range collection {
		_, execErr := tx.Exec(context.Background(), insertSql, metrics.ID, metrics.MType, metrics.Delta, metrics.Value, metrics.Hash)
		if execErr != nil {
			fmt.Printf("Ошибка при восстановлении метрики (id: %s, mtype: %s): %v\n", metrics.ID, metrics.MType, execErr)
			txErr = execErr
			return
		}
	}
}

func (pg *pg) Ping(ctx context.Context) error {
	return pg.dbpool.Ping(ctx)
}
