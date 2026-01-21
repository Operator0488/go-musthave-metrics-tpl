package repository

import (
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"log"
	"strconv"
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
	storage map[string]models.Metrics

	mu sync.RWMutex
}

func NewMaps() *Maps {
	return &Maps{
		storage: make(map[string]models.Metrics),
	}
}

// SaverValue -
func (m *Maps) SaverValue(name string, value float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Gauge {
			return models.ErrorDiffType
		}
		log.Println(value)
		*metric.Value = value
		return nil
	}

	var metric models.Metrics
	metric.MType = models.Gauge
	metric.Value = &value

	metric.ID = strconv.Itoa(len(m.storage) + 1)
	m.storage[name] = metric

	return nil
}

// IncrementValue -
func (m *Maps) IncrementValue(name string, value int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Counter {
			return models.ErrorGetValue
		}
		*metric.Delta += value
		return nil
	}

	var metric models.Metrics
	metric.MType = models.Counter
	metric.Delta = &value

	metric.ID = strconv.Itoa(len(m.storage) + 1)
	m.storage[name] = metric

	return nil
}

// GetValueGauge -
func (m *Maps) GetValueGauge(s string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[s]; ok {
		if metric.Value != nil {
			n := metric.Value
			return *n, nil
		} else {
			return 0, models.ErrorDiffType
		}
	}

	return 0, models.ErrorNotDB
}

// GetValueCounter -
func (m *Maps) GetValueCounter(s string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[s]; ok {
		if metric.Delta != nil {
			n := metric.Delta
			return *n, nil
		} else {
			return 0, models.ErrorDiffType
		}

	}

	return 0, models.ErrorNotDB
}

// GetValues -
func (m *Maps) GetValues() (map[string]any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	storage := make(map[string]interface{})

	for k, v := range m.storage {
		switch v.MType {
		case models.Gauge:
			value := v.Value
			storage[k] = *value
		case models.Counter:
			value := v.Delta
			storage[k] = *value
		default:
			return nil, models.ErrorUnType
		}
	}

	return storage, nil
}
