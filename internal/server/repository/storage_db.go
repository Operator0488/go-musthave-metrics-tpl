package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"sync"
)

type PostgresStorage struct {
	db *sql.DB
	mu sync.RWMutex

	log logger.Logger
}

func NewPostgresStorage(log logger.Logger, db *sql.DB) (*PostgresStorage, error) {
	var store = &PostgresStorage{
		log: log,
		db:  db,
	}

	return store, nil
}

func (m *PostgresStorage) SaveValue(ctx context.Context, s string, f float64) error {
	//noinspection SqlResolve
	query := `INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2) ON CONFLICT (id) 
    DO UPDATE SET 
                value = EXCLUDED.value, delta = NULL
    			WHERE metrics.type = 'gauge';`

	res, err := m.db.ExecContext(ctx, query, s, f)
	if err != nil {
		return fmt.Errorf("ошибка сохранения значения gauge: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("ошибка формата %s не gauge", s)
	}

	return nil
}

func (m *PostgresStorage) IncrementValue(ctx context.Context, s string, i int64) error {
	//noinspection SqlResolve
	query := `INSERT INTO metrics (id, type, delta) VALUES ($1, 'counter', $2) ON CONFLICT (id) 
    DO UPDATE SET 
                delta = metrics.delta + EXCLUDED.delta, value = NULL
    			WHERE metrics.type = 'counter';`

	res, err := m.db.ExecContext(ctx, query, s, i)
	if err != nil {
		return fmt.Errorf("ошибка сохранения значения counter: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("ошибка формата %s не counter", s)
	}

	return nil
}

func (m *PostgresStorage) GetValueGauge(ctx context.Context, s string) (float64, error) {
	//noinspection SqlResolve
	query := `SELECT type, value FROM metrics WHERE id=$1;`

	var t string
	var f float64
	err := m.db.QueryRowContext(ctx, query, s).Scan(&t, &f)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения значения gauge: %w", err)
	}

	if t != models.Gauge {
		return 0, fmt.Errorf("ошибка получения значения gauge: %w", err)
	}

	return f, nil
}

func (m *PostgresStorage) GetValueCounter(ctx context.Context, s string) (int64, error) {
	//noinspection SqlResolve
	query := `SELECT type, delta FROM metrics WHERE id=$1;`

	var t string
	var f int64
	err := m.db.QueryRowContext(ctx, query, s).Scan(&t, &f)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения значения counter: %w", err)
	}

	if t != models.Counter {
		return 0, fmt.Errorf("ошибка получения значения counter: %w", err)
	}

	return f, nil
}

func (m *PostgresStorage) GetValues(ctx context.Context) (map[string]any, error) {
	//noinspection SqlResolve
	query := `SELECT id, type, value, delta FROM metrics`

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения значений %w", err)
	}

	storage := make(map[string]interface{})

	for rows.Next() {
		var id, t string
		var value sql.NullFloat64
		var delta sql.NullInt64
		err = rows.Scan(&id, &t, &value, &delta)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения значений %w", err)
		}
		switch t {
		case models.Gauge:
			storage[id] = value.Float64
		case models.Counter:
			storage[id] = delta.Int64
		default:
			return nil, models.ErrorUnType
		}
	}

	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения значений %w", err)
	}

	return storage, nil
}

func (m *PostgresStorage) PingDB(ctx context.Context) error {
	err := m.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("ошибка ping бд: %w", err)
	}
	return nil
}

func (m *PostgresStorage) SaveValues(ctx context.Context, reqs []models.PostUpdateRequest) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	for _, req := range reqs {
		switch req.MType {
		case models.Gauge:
			//noinspection SqlResolve
			query := `INSERT INTO metrics (id, type, value) VALUES ($1, 'gauge', $2) ON CONFLICT (id) 
    DO UPDATE SET 
                value = EXCLUDED.value, delta = NULL
    			WHERE metrics.type = 'gauge';`

			_, err := tx.ExecContext(ctx,
				query, req.ID, req.Value)

			if err != nil {
				var errs error
				er := tx.Rollback()
				if er != nil {
					errs = errors.Join(fmt.Errorf("ошибка начала транзакции: %w", err), fmt.Errorf("ошибка отмены транзакции: %w", er))
				} else {
					errs = fmt.Errorf("ошибка начала транзакции: %w", err)
				}
				return errs
			}
		case models.Counter:
			//noinspection SqlResolve
			query := `INSERT INTO metrics (id, type, delta) VALUES ($1, 'counter', $2) ON CONFLICT (id) 
    DO UPDATE SET 
                delta = metrics.delta + EXCLUDED.delta, value = NULL
    			WHERE metrics.type = 'counter';`

			_, err := tx.ExecContext(ctx,
				query, req.ID, req.Delta)

			if err != nil {
				var errs error
				er := tx.Rollback()
				if er != nil {
					errs = errors.Join(fmt.Errorf("ошибка начала транзакции: %w", err), fmt.Errorf("ошибка отмены транзакции: %w", er))
				} else {
					errs = fmt.Errorf("ошибка начала транзакции: %w", err)
				}
				return errs
			}
		}

	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ошибка при коммите: %w", err)
	}

	return nil
}
