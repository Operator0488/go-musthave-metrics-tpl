package service

import (
	"context"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"reflect"
	"sync"
	"testing"
)

func TestNewMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*model.Stat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStatsManager(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *StatsManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatsManager(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatsManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsManager_GetMap(t *testing.T) {
	type fields struct {
		ctx context.Context
		m   map[string]*model.Stat
		mu  sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*model.Stat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &StatsManager{
				ctx: tt.fields.ctx,
				m:   tt.fields.m,
				mu:  tt.fields.mu,
			}
			if got := m.GetMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatsManager_WriteStats(t *testing.T) {
	type fields struct {
		ctx context.Context
		m   map[string]*model.Stat
		mu  sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &StatsManager{
				ctx: tt.fields.ctx,
				m:   tt.fields.m,
				mu:  tt.fields.mu,
			}
			m.WriteStats()
		})
	}
}

func TestStatsManager_GetMap_ReturnsSnapshot(t *testing.T) {
	mn := NewStatsManager(context.Background())
	mn.WriteStats()

	// Снимок №1
	snap1 := mn.GetMap()

	// 1) В snapshot должны быть ключи
	if _, ok := snap1[model.Alloc]; !ok {
		t.Fatalf("snapshot missing key %q", model.Alloc)
	}
	if _, ok := snap1[model.PollCount]; !ok {
		t.Fatalf("snapshot missing key %q", model.PollCount)
	}

	// 2) Значения должны совпадать с новым снимком (в момент получения)
	alloc1 := snap1[model.Alloc].Value
	poll1 := snap1[model.PollCount].Value

	snap2 := mn.GetMap()
	if snap2[model.Alloc].Value != alloc1 {
		t.Fatalf("Alloc mismatch between snapshots: snap1=%v snap2=%v", alloc1, snap2[model.Alloc].Value)
	}
	if snap2[model.PollCount].Value != poll1 {
		t.Fatalf("PollCount mismatch between snapshots: snap1=%v snap2=%v", poll1, snap2[model.PollCount].Value)
	}

	// 3) Указатели на Stat должны быть разными (иначе это не snapshot, а утечка внутренностей)
	if snap1[model.Alloc] == snap2[model.Alloc] {
		t.Fatal("expected different *Stat pointers for Alloc between snapshots, got same pointer")
	}

	// 4) Изменения в snapshot НЕ должны протекать внутрь менеджера

	// 4.1) Меняем значение внутри Stat в snapshot
	snap1[model.Alloc].Value = 123456789

	after := mn.GetMap()
	if after[model.Alloc].Value == 123456789 {
		t.Fatal("snapshot mutation leaked into manager state (Alloc.Value changed)")
	}

	// 4.2) Подменяем указатель в snapshot (вообще не должно влиять на менеджер)
	snap1[model.Alloc] = &model.Stat{Type: "gauge", Value: 42}

	after2 := mn.GetMap()
	if after2[model.Alloc].Value == 42 {
		t.Fatal("snapshot map pointer replacement leaked into manager state")
	}
}
