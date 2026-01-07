package repository

import (
	"context"
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"strconv"
	"sync"
)

type MemStorage interface {
	SaverValue(string, float64) error
	IncrementValue(string, int64) error
}

type Maps struct {
	ctx     context.Context
	storage map[string]models.Metrics

	mu sync.RWMutex
}

func NewMaps(ctx context.Context) *Maps {
	return &Maps{
		ctx:     ctx,
		storage: make(map[string]models.Metrics),
	}
}

// SaverValue -
func (m *Maps) SaverValue(name string, value float64) error {

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Gauge {
			return fmt.Errorf("Ошибка: переменная в базе другого типа")
		}
		metric.Value = &value
		return nil
	}

	var metric models.Metrics
	metric.MType = models.Gauge
	metric.Value = &value

	m.mu.Lock()
	defer m.mu.Unlock()

	metric.ID = strconv.Itoa(len(m.storage) + 1)
	m.storage[name] = metric

	return nil
}

// IncrementValue -
func (m *Maps) IncrementValue(name string, value int64) error {

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Counter {
			return fmt.Errorf("Ошибка: переменная в базе другого типа")
		}
		*metric.Delta += value
		return nil
	}

	var metric models.Metrics
	metric.MType = models.Counter
	metric.Delta = &value

	m.mu.Lock()
	defer m.mu.Unlock()

	metric.ID = strconv.Itoa(len(m.storage) + 1)
	m.storage[name] = metric

	return nil
}
