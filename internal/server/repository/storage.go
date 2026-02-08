package repository

import (
	"encoding/json"
	"errors"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type MemStorage interface {
	SaverValue(string, float64) error
	IncrementValue(string, int64) error
	GetValueGauge(string) (float64, error)
	GetValueCounter(string) (int64, error)
	GetValues() (map[string]any, error)
	SaveData(models.MetricStore) error
	GetData() error
	AddData(models.MetricStore) error
	Snapshot() error
}

type Maps struct {
	storage map[string]models.MetricStore
	file    *os.File

	mu         sync.RWMutex
	syncRecord bool
}

func loadFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func NewMaps(checkInit bool, path string, interval int) (*Maps, error) {
	store := make(map[string]models.MetricStore)

	file, err := loadFile(path)
	if err != nil {
		return nil, err
	}

	var maps = &Maps{
		storage:    store,
		file:       file,
		syncRecord: interval == 0,
	}

	if checkInit {
		err = maps.GetData()
		if err != nil {
			return nil, err
		}
	}

	return maps, nil
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

	if m.syncRecord {
		err := m.SaveData(metric)
		return err
	}

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

	if m.syncRecord {
		err := m.SaveData(metric)
		return err
	}

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

func (m *Maps) SaveData(request models.MetricStore) error {
	enc := json.NewEncoder(m.file)
	enc.SetIndent("", "\t")

	err := enc.Encode(request)
	if err != nil {
		return err
	}

	return nil
}

func (m *Maps) GetData() error {
	dec := json.NewDecoder(m.file)

	for {
		var ms models.MetricStore
		err := dec.Decode(&ms)
		if errors.Is(err, io.EOF) {
			break
		}
		err = m.AddData(ms)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Maps) AddData(ms models.MetricStore) error {
	if ms.MType == models.Gauge {
		err := m.SaverValue(ms.ID, ms.Value)
		if err != nil {
			return err
		}
	}

	if ms.MType == models.Counter {
		err := m.IncrementValue(ms.ID, ms.Delta)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Maps) Snapshot() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	err := m.file.Truncate(0)
	if err != nil {
		return err
	}
	for _, v := range m.storage {
		err := m.SaveData(v)
		if err != nil {
			return err
		}
	}
	return nil
}
