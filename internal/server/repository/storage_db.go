package repository

import (
	"context"
	"database/sql"
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

// noinspection SqlResolve
const (
	qInsertGauge = `
INSERT INTO metrics (id, type, value)
VALUES ($1, 'gauge', $2)
ON CONFLICT (id) DO UPDATE
SET value = EXCLUDED.value, delta = NULL
WHERE metrics.type = 'gauge';
`
	qInsertCounter = `
INSERT INTO metrics (id, type, delta)
VALUES ($1, 'counter', $2)
ON CONFLICT (id) DO UPDATE
SET delta = metrics.delta + EXCLUDED.delta, value = NULL
WHERE metrics.type = 'counter';
`
	qSelectCounter = `
SELECT type, delta FROM metrics WHERE id=$1;
`
	qSelectGauge = `
SELECT type, value FROM metrics WHERE id=$1;
`
	qSelectAll = `
SELECT id, type, value, delta FROM metrics
`
)

func (m *PostgresStorage) SaveValue(ctx context.Context, s string, f float64) error {
	res, err := m.db.ExecContext(ctx, qInsertGauge, s, f)
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
	res, err := m.db.ExecContext(ctx, qInsertCounter, s, i)
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
	var t string
	var f float64
	err := m.db.QueryRowContext(ctx, qSelectGauge, s).Scan(&t, &f)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения значения gauge: %w", err)
	}

	if t != models.Gauge {
		return 0, fmt.Errorf("ошибка получения значения gauge: %w", err)
	}

	return f, nil
}

func (m *PostgresStorage) GetValueCounter(ctx context.Context, s string) (int64, error) {
	var t string
	var f int64
	err := m.db.QueryRowContext(ctx, qSelectCounter, s).Scan(&t, &f)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения значения counter: %w", err)
	}

	if t != models.Counter {
		return 0, fmt.Errorf("ошибка получения значения counter: %w", err)
	}

	return f, nil
}

func (m *PostgresStorage) GetValues(ctx context.Context) (map[string]any, error) {
	rows, err := m.db.QueryContext(ctx, qSelectAll)
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

func (m *PostgresStorage) SaveValues(ctx context.Context, reqs []models.PostUpdateRequest) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("ошибка коммита: %w", err)
		}
	}()

	for _, req := range reqs {
		switch req.MType {
		case models.Gauge:

			if req.Value == nil {
				return fmt.Errorf("ошибка записи gauge %q: нет значения", req.ID)
			}
			res, e := tx.ExecContext(ctx,
				qInsertGauge, req.ID, req.Value)

			if e != nil {
				return fmt.Errorf("ошибка записи gauge %q: %w", req.ID, e)
			}

			rows, _ := res.RowsAffected()
			if rows == 0 {
				return fmt.Errorf("%w: %q ожидали gauge", models.ErrorUnType, req.ID)
			}

		case models.Counter:

			if req.Delta == nil {
				return fmt.Errorf("ошибка записи counter %q: нет значения", req.ID)
			}

			res, e := tx.ExecContext(ctx,
				qInsertCounter, req.ID, req.Delta)

			if e != nil {
				return fmt.Errorf("ошибка записи counter %q: %w", req.ID, e)
			}

			rows, _ := res.RowsAffected()
			if rows == 0 {
				return fmt.Errorf("%w: %q ожидали сounter", models.ErrorUnType, req.ID)
			}

		default:
			return models.ErrorUnType
		}

	}

	return nil
}
