package service

import (
	"context"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository/mock"
	"reflect"
	"testing"
)

func TestNewStorageService(t *testing.T) {
	type args struct {
		ctx     context.Context
		storage repository.MemStorage
	}
	tests := []struct {
		name string
		args args
		want *StorageService
	}{
		{
			name: "test new storage, type",
			args: args{
				ctx:     context.Background(),
				storage: &repository.Maps{},
			},
			want: &StorageService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorageService(tt.args.ctx, tt.args.storage); !reflect.DeepEqual(reflect.TypeOf(got).Name(), reflect.TypeOf(tt.want).Name()) {
				t.Errorf("NewStorageService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageService_Sender(t *testing.T) {
	type fields struct {
		ctx     context.Context
		storage repository.MemStorage
	}
	type args struct {
		req *models.PostUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test storage service #1",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Counter,
					Name:  "Count",
					Value: "123",
				},
			},
			wantErr: false,
		},
		{
			name: "test storage service #2",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Gauge,
					Name:  "Counter",
					Value: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "test storage service #3",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Counter,
					Name:  "Counter",
					Value: "fsgfsffg",
				},
			},
			wantErr: true,
		},
		{
			name: "test storage service #4",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Counter,
					Name:  "Gauge",
					Value: "123.01",
				},
			},
			wantErr: true,
		},
		{
			name: "test storage service #5",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Gauge,
					Name:  "Gauge",
					Value: "123.01",
				},
			},
			wantErr: false,
		},
		{
			name: "test storage service #6",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Gauge,
					Name:  "Counter",
					Value: "123.01",
				},
			},
			wantErr: true,
		},
		{
			name: "test storage service #7",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Gauge,
					Name:  "Gauge",
					Value: "123",
				},
			},
			wantErr: false,
		},
		{
			name: "test storage service #8",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  "node",
					Name:  "Counter",
					Value: "123.01",
				},
			},
			wantErr: true,
		},
		{
			name: "test storage service #9",
			fields: fields{
				ctx:     context.Background(),
				storage: mock.NewMapsMock(),
			},
			args: args{
				req: &models.PostUpdateRequest{
					Type:  models.Gauge,
					Name:  "Counter",
					Value: "fgjn",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageService{
				ctx:     tt.fields.ctx,
				storage: tt.fields.storage,
			}
			if err := s.SenderPostUpdate(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("SenderPostUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
