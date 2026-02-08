package service

import (
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"math/rand/v2"
	"runtime"
	"sync"
)

type StatsManager struct {
	m  map[string]*model.Stat
	mu sync.RWMutex
}

type Manager interface {
	WriteStats()
	GetMap() map[string]*model.Stat
}

// NewStatsManager -
func NewStatsManager() *StatsManager {
	var s = StatsManager{
		m: NewMap(),
	}
	s.WriteStats()

	return &s
}

func NewMap() map[string]*model.Stat {
	return map[string]*model.Stat{
		model.GCCPUFraction: {"gauge", 0.0},
		model.NumForcedGC:   {"gauge", 0.0},
		model.NumGC:         {"gauge", 0.0},
		model.PauseTotalNs:  {"gauge", 0.0},
		model.LastGC:        {"gauge", 0.0},
		model.MCacheSys:     {"gauge", 0.0},
		model.MCacheInuse:   {"gauge", 0.0},
		model.NextGC:        {"gauge", 0.0},
		model.OtherSys:      {"gauge", 0.0},
		model.GCSys:         {"gauge", 0.0},
		model.BuckHashSys:   {"gauge", 0.0},
		model.MSpanSys:      {"gauge", 0.0},
		model.MSpanInuse:    {"gauge", 0.0},
		model.StackSys:      {"gauge", 0.0},
		model.StackInuse:    {"gauge", 0.0},
		model.HeapObjects:   {"gauge", 0.0},
		model.HeapReleased:  {"gauge", 0.0},
		model.HeapInuse:     {"gauge", 0.0},
		model.HeapIdle:      {"gauge", 0.0},
		model.HeapSys:       {"gauge", 0.0},
		model.HeapAlloc:     {"gauge", 0.0},
		model.Frees:         {"gauge", 0.0},
		model.Mallocs:       {"gauge", 0.0},
		model.Lookups:       {"gauge", 0.0},
		model.Sys:           {"gauge", 0.0},
		model.TotalAlloc:    {"gauge", 0.0},
		model.Alloc:         {"gauge", 0.0},
		model.RandomValue:   {"gauge", 0.0},
		model.PollCount:     {"counter", 0.0},
	}
}

func (m *StatsManager) WriteStats() {
	var s runtime.MemStats
	runtime.SetMutexProfileFraction(10)
	runtime.ReadMemStats(&s)

	m.mu.Lock()
	m.m[model.Alloc].Value = float64(s.Alloc)
	m.m[model.Frees].Value = float64(s.Frees)
	m.m[model.HeapAlloc].Value = float64(s.HeapAlloc)
	m.m[model.Sys].Value = float64(s.Sys)
	m.m[model.BuckHashSys].Value = float64(s.BuckHashSys)
	m.m[model.GCSys].Value = float64(s.GCSys)
	m.m[model.GCCPUFraction].Value = s.GCCPUFraction
	m.m[model.HeapIdle].Value = float64(s.HeapIdle)
	m.m[model.HeapObjects].Value = float64(s.HeapObjects)
	m.m[model.HeapInuse].Value = float64(s.HeapInuse)
	m.m[model.HeapReleased].Value = float64(s.HeapReleased)
	m.m[model.HeapSys].Value = float64(s.HeapSys)
	m.m[model.LastGC].Value = float64(s.LastGC)
	m.m[model.Lookups].Value = float64(s.Lookups)
	m.m[model.MCacheInuse].Value = float64(s.MCacheInuse)
	m.m[model.MCacheSys].Value = float64(s.MCacheSys)
	m.m[model.MSpanInuse].Value = float64(s.MSpanInuse)
	m.m[model.MSpanSys].Value = float64(s.MSpanSys)
	m.m[model.Mallocs].Value = float64(s.Mallocs)
	m.m[model.NextGC].Value = float64(s.NextGC)
	m.m[model.NumForcedGC].Value = float64(s.NumForcedGC)
	m.m[model.NumGC].Value = float64(s.NumGC)
	m.m[model.OtherSys].Value = float64(s.OtherSys)
	m.m[model.PauseTotalNs].Value = float64(s.PauseTotalNs)
	m.m[model.StackInuse].Value = float64(s.StackInuse)
	m.m[model.StackSys].Value = float64(s.StackSys)
	m.m[model.TotalAlloc].Value = float64(s.TotalAlloc)
	m.m[model.PollCount].Value = 1
	m.m[model.RandomValue].Value = rand.Float64()
	m.mu.Unlock()

}

func (m *StatsManager) GetMap() map[string]*model.Stat {
	m.mu.Lock()
	defer m.mu.Unlock()

	cp := make(map[string]*model.Stat, len(m.m))
	for k, v := range m.m {
		if v == nil {
			cp[k] = nil
			continue
		}
		s := *v
		cp[k] = &s
	}
	return cp
}
