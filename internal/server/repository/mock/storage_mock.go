package mock

import (
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"strconv"
)

type MapsMock struct {
	storage map[string]models.Metrics
}

func NewMapsMock() *MapsMock {
	delta := int64(10)
	value := 11.00
	m := MapsMock{
		storage: map[string]models.Metrics{
			"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Delta: &delta},
			"Counter": models.Metrics{ID: "1", MType: models.Counter, Value: &value},
		},
	}
	return &m
}

func (m *MapsMock) SaverValue(name string, value float64) error {

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

	metric.ID = strconv.Itoa(len(m.storage) + 1)
	m.storage[name] = metric

	return nil
}

func (m *MapsMock) IncrementValue(name string, value int64) error {

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

func (m *MapsMock) PrintAll() {

}
