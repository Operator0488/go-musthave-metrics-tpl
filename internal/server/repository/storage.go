package repository

import (
	"context"
	"fmt"
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
	GetValues() (map[string]interface{}, error)
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
	m.mu.Lock()
	defer m.mu.Unlock()

	if metric, ok := m.storage[name]; ok {
		if metric.MType != models.Gauge {
			return fmt.Errorf("Ошибка: переменная в базе другого типа")
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
			return fmt.Errorf("Ошибка: переменная в базе другого типа")
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
			return 0, fmt.Errorf("Ошибка: переменная в базе другого типа")
		}
	}

	return 0, fmt.Errorf("Ошибка: переменной нет в базе")
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
			return 0, fmt.Errorf("Ошибка: переменная в базе другого типа")
		}

	}

	return 0, fmt.Errorf("Ошибка: переменной нет в базе")
}

// GetValues -
func (m *Maps) GetValues() (map[string]interface{}, error) {
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
			return nil, fmt.Errorf("Есть непредвиденный тип данных в базе")
		}
	}

	return storage, nil
}
