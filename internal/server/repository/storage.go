package repository

import (
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"sync"
)

type MemStorage interface {
	SaverValue(string, float64) error
	IncrementValue(string, int64) error
	GetValueGauge(string) (float64, error)
	GetValueCounter(string) (int64, error)
	GetValues() (map[string]any, error)
}

type Maps struct {
	storage map[string]models.MetricStore

	mu sync.RWMutex
}

func NewMaps() *Maps {
	return &Maps{
		storage: make(map[string]models.MetricStore),
	}
}

func (m *Maps) SaverValue(name string, value float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Gauge {
			return models.ErrorDiffType
		}
		metric.Value = value
		m.storage[name] = metric
		return nil
	}

	var metric models.MetricStore
	metric.MType = models.Gauge
	metric.Value = value

	metric.ID = name
	m.storage[name] = metric

	return nil
}

func (m *Maps) IncrementValue(name string, value int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Counter {
			return models.ErrorGetValue
		}
		metric.Delta += value
		m.storage[name] = metric
		return nil
	}

	var metric models.MetricStore
	metric.MType = models.Counter
	metric.Delta = value

	metric.ID = name
	m.storage[name] = metric

	return nil
}

func (m *Maps) GetValueGauge(name string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if metric, ok := m.storage[name]; ok {
		return metric.Value, nil
	}

	return 0, models.ErrorNotDB
}

func (m *Maps) GetValueCounter(name string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if metric, ok := m.storage[name]; ok {
		return metric.Delta, nil
	}

	return 0, models.ErrorNotDB
}

func (m *Maps) GetValues() (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	storage := make(map[string]interface{})

	for k, v := range m.storage {
		switch v.MType {
		case models.Gauge:
			storage[k] = v.Value
		case models.Counter:
			storage[k] = v.Delta
		default:
			return nil, models.ErrorUnType
		}
	}

	return storage, nil
}
