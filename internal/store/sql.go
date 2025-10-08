package store

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/faultTolerance/pgretry"
	"github.com/dontagr/metric/internal/server/service/interfaces"
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
	insertSQL = `INSERT INTO metric (id, mtype, delta, value, hash) 
	  VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, mtype) DO UPDATE SET
	  mtype = EXCLUDED.mtype,
	  delta = EXCLUDED.delta,
	  value = EXCLUDED.value,
	  hash = EXCLUDED.hash`
	bulkInsertSQL = `INSERT INTO metric (id, mtype, delta, value, hash)
	  VALUES %s ON CONFLICT (id, mtype) DO UPDATE SET
	  delta = EXCLUDED.delta,
	  value = EXCLUDED.value,
	  hash = EXCLUDED.hash`
	searchSQL = `SELECT id, mtype, delta, value, hash 
      FROM metric 
      WHERE id=$1 AND mtype=$2`
	selectAllSQL = `SELECT id, mtype, delta, value, hash FROM metric`
	truncateSQL  = `TRUNCATE TABLE metric`
)

type pg struct {
	dbpool *pgretry.PgxRetry
	log    *zap.SugaredLogger
	name   string
	mx     sync.RWMutex
}

func RegisterStorePG(log *zap.SugaredLogger, ms interfaces.IStoreFactory, dbpool *pgretry.PgxRetry, lc fx.Lifecycle) {
	pg := pg{
		dbpool: dbpool,
		name:   models.StorePg,
		log:    log,
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

func (pg *pg) LoadMetric(id string, mType string) (*models.Metrics, error) {
	pg.mx.RLock()
	defer pg.mx.RUnlock()

	var metrics models.Metrics
	err := pg.dbpool.QueryRow(context.Background(), searchSQL, id, mType).Scan(
		&metrics.ID,
		&metrics.MType,
		&metrics.Delta,
		&metrics.Value,
		&metrics.Hash,
	)
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

func (pg *pg) SaveMetric(metrics *models.Metrics) error {
	pg.mx.Lock()
	defer pg.mx.Unlock()

	id, mtype, delta, value, hash := pg.unpack(metrics)
	_, err := pg.dbpool.Exec(context.Background(), insertSQL, id, mtype, delta, value, hash)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении метрики: %w", err)
	}

	return nil
}

func (pg *pg) BulkSaveMetric(metrics map[string]*models.Metrics) error {
	pg.mx.Lock()
	tx, txErr := pg.dbpool.Begin(context.Background())
	if txErr != nil {
		pg.mx.Unlock()
		return fmt.Errorf("ошибка начала транзакции: %w", txErr)
	}
	defer func(txErr *error) {
		if *txErr != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				pg.log.Errorf("ошибка отката транзакции: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(context.Background()); commitErr != nil {
				pg.log.Errorf("ошибка при коммите транзакции: %v", commitErr)
			}
		}
		pg.mx.Unlock()
	}(&txErr)

	values := make([]interface{}, 0, len(metrics)*5)
	valueStrings := make([]string, 0, len(metrics))
	i := 0

	for _, metric := range metrics {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i+1, i+2, i+3, i+4, i+5))
		values = append(values, metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)
		i += 5
	}

	sqlStr := fmt.Sprintf(bulkInsertSQL, strings.Join(valueStrings, ","))
	_, execErr := tx.Exec(context.Background(), sqlStr, values...)
	if execErr != nil {
		txErr = execErr
		return fmt.Errorf("ошибка при массовом обновлении метрик: %w", execErr)
	}

	return nil
}

func (pg *pg) unpack(metrics *models.Metrics) (string, string, *int64, *float64, string) {
	return metrics.ID, metrics.MType, metrics.Delta, metrics.Value, metrics.Hash
}

func (pg *pg) ListMetric() (map[string]*models.Metrics, error) {
	pg.mx.RLock()
	defer pg.mx.RUnlock()

	r := make(map[string]*models.Metrics)
	rows, err := pg.dbpool.Query(context.Background(), selectAllSQL)
	if err != nil {
		return r, fmt.Errorf("ошибка при извлечении метрик: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var metrics models.Metrics
		err := rows.Scan(&metrics.ID, &metrics.MType, &metrics.Delta, &metrics.Value, &metrics.Hash)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании метрики: %w", err)
		}
		r[fmt.Sprintf("%s_%s", metrics.MType, metrics.ID)] = &metrics
	}

	return r, nil
}

func (pg *pg) RestoreMetricCollection(ctx context.Context, collection map[string]*models.Metrics) error {
	pg.mx.Lock()
	tx, txErr := pg.dbpool.Begin(ctx)
	if txErr != nil {
		pg.mx.Unlock()
		return fmt.Errorf("ошибка начала транзакции: %w", txErr)
	}
	defer func(txErr *error) {
		if *txErr != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				pg.log.Errorf("ошибка отката транзакции: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(context.Background()); commitErr != nil {
				pg.log.Errorf("ошибка при коммите транзакции: %v", commitErr)
			}
		}
		pg.mx.Unlock()
	}(&txErr)

	_, execErr := tx.Exec(context.Background(), truncateSQL)
	if execErr != nil {
		txErr = execErr
		return fmt.Errorf("ошибка при восстановлении метрики: %w", execErr)
	}

	values := make([]interface{}, 0, len(collection)*5)
	valueStrings := make([]string, 0, len(collection))
	i := 0

	for _, metrics := range collection {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i+1, i+2, i+3, i+4, i+5))
		values = append(values, metrics.ID, metrics.MType, metrics.Delta, metrics.Value, metrics.Hash)
		i += 5
	}

	sqlStr := fmt.Sprintf(bulkInsertSQL, strings.Join(valueStrings, ","))
	_, execErr = tx.Exec(context.Background(), sqlStr, values...)
	if execErr != nil {
		txErr = execErr
		return fmt.Errorf("ошибка при восстановлении метрик: %w", execErr)
	}

	return nil
}

func (pg *pg) Ping(ctx context.Context) error {
	return pg.dbpool.Ping(ctx)
}
