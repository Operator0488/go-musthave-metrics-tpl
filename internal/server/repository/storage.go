package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:generate mockgen -source=storage.go -destination=./mocks/mock_storage.go -package=mocks
type MemStorage interface {
	SaverValue(string, float64) error
	IncrementValue(string, int64) error
	GetValueGauge(string) (float64, error)
	GetValueCounter(string) (int64, error)
	GetValues() (map[string]any, error)
	GetData() error
	AddData(models.MetricStore) error
}

type Maps struct {
	storage map[string]models.MetricStore
	file    *os.File

	mu         sync.RWMutex
	muFile     sync.RWMutex
	syncRecord bool

	log logger.Logger
}

func loadFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("ошибка создания директории: %w", err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания файла: %w", err)
	}

	return file, nil
}

func NewMaps(log logger.Logger, checkInit bool, path string, interval int) (*Maps, error) {
	store := make(map[string]models.MetricStore)

	file, err := loadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки файла: %w", err)
	}

	var maps = &Maps{
		storage:    store,
		file:       file,
		syncRecord: interval == 0,
		log:        log,
	}

	if checkInit {
		err = maps.GetData()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных: %w", err)
		}
	}

	maps.snapshot(interval)

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
		if m.syncRecord {
			err := m.saveData(metric)
			return err
		}
		return nil
	}

	var metric models.MetricStore
	metric.MType = models.Gauge
	metric.Value = value

	metric.ID = name
	m.storage[name] = metric

	if m.syncRecord {
		err := m.saveData(metric)
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
		if m.syncRecord {
			metric.Delta = value
			err := m.saveData(metric)
			return err
		}
		return nil
	}

	var metric models.MetricStore
	metric.MType = models.Counter
	metric.Delta = value

	metric.ID = name
	m.storage[name] = metric

	if m.syncRecord {
		err := m.saveData(metric)
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

func (m *Maps) saveData(request models.MetricStore) error {
	m.muFile.Lock()
	defer m.muFile.Unlock()

	enc := json.NewEncoder(m.file)
	enc.SetIndent("", "\t")

	err := enc.Encode(request)
	if err != nil {
		return fmt.Errorf("ошибка сохранения models.MetricStore: %w", err)
	}

	return nil
}

func (m *Maps) GetData() error {
	dec := json.NewDecoder(m.file)

	for {
		var ms models.MetricStore
		err := dec.Decode(&ms)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, io.ErrUnexpectedEOF) {
				m.log.Info("Записали не все данные", zap.Error(err))
				break
			}
			return fmt.Errorf("ошибка получения данных: %w", err)
		}
		err = m.AddData(ms)
		if err != nil {
			return fmt.Errorf("ошибка добавление данных в storage: %w", err)
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

func (m *Maps) snapshot(interval int) {
	if interval <= 0 {
		return
	}

	tick := time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		defer tick.Stop()

		for range tick.C {

			m.muFile.Lock()
			err := m.file.Truncate(0)
			if err != nil {
				m.muFile.Unlock()
				m.log.Info("ошибка очистки файла",
					zap.Error(err))
				break
			}
			_, err = m.file.Seek(0, 0)
			if err != nil {
				m.muFile.Unlock()
				m.log.Info("ошибка сброса записи в начало",
					zap.Error(err))
				break
			}
			m.muFile.Unlock()

			m.mu.RLock()
			for _, v := range m.storage {
				err = m.saveData(v)
				if err != nil {
					m.mu.RUnlock()
					m.log.Info("ошибка сохранения данных",
						zap.Error(err))
				}
			}
			m.mu.RUnlock()
		}
	}()
}
