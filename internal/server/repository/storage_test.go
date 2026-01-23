package repository

import (
	"context"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"reflect"
	"sync"
	"testing"
)

func TestMaps_IncrementValue(t *testing.T) {
	type fields struct {
		ctx     context.Context
		storage map[string]models.Metrics
		mu      sync.RWMutex
	}
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test storage incrementvalues #1",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Counter",
				value: 12,
			},
			wantErr: false,
		},
		{
			name: "test storage incrementvalues #2",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Counter",
				value: 12,
			},
			wantErr: false,
		},
		{
			name: "test storage incrementvalues #3",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Type",
				value: 14,
			},
			wantErr: false,
		},
		{
			name: "test storage incrementvalues #4",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Gauge",
				value: 12,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Maps{
				storage: tt.fields.storage,
				mu:      tt.fields.mu,
			}
			if err := m.IncrementValue(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("IncrementValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "test storage incrementvalues #2" {
				val := *tt.fields.storage["Counter"].Delta
				if val != 23 {
					t.Errorf("IncrementValue() val = %v, wantErr %v", val, 23)
				}
			}
			if tt.name == "test storage incrementvalues #3" {
				val := *tt.fields.storage["Type"].Delta
				if val != 14 {
					t.Errorf("IncrementValue() val = %v, wantErr %v", val, 14)
				}
			}
		})
	}
}

func TestMaps_SaverValue(t *testing.T) {
	type fields struct {
		ctx     context.Context
		storage map[string]models.Metrics
		mu      sync.RWMutex
	}
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test storage savervalues #1",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Gauge",
				value: 12.0,
			},
			wantErr: false,
		},
		{
			name: "test storage savervalues #2",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Gauge",
				value: 12,
			},
			wantErr: false,
		},
		{
			name: "test storage savervalues #3",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Counter",
				value: 12,
			},
			wantErr: true,
		},
		{
			name: "test storage savervalues #4",
			fields: fields{
				ctx: context.Background(),
				storage: map[string]models.Metrics{
					"Gauge":   models.Metrics{ID: "1", MType: models.Gauge, Value: PtrFloat64(10.012)},
					"Counter": models.Metrics{ID: "1", MType: models.Counter, Delta: PtrInt64(11.00)},
				},
			},
			args: args{
				name:  "Gauge2",
				value: 12.000000000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Maps{
				storage: tt.fields.storage,
				mu:      tt.fields.mu,
			}
			if err := m.SaverValue(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SaverValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMaps(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *Maps
	}{
		{
			name: "test storage newmaps #1",
			args: args{
				ctx: context.Background(),
			},
			want: &Maps{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMaps(); !reflect.DeepEqual(reflect.TypeOf(got).Name(), reflect.TypeOf(tt.want).Name()) {
				t.Errorf("NewMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func PtrInt64(v int64) *int64 {
	return &v
}

func PtrFloat64(v float64) *float64 {
	return &v
}
