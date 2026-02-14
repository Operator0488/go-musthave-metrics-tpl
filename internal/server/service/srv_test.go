package service

import (
	"errors"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	mocklog "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger/mocks"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestNewStorageService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		storage repository.MemStorage
		log     logger.Logger
	}

	tests := []struct {
		name string
		args args
		want *StorageService
	}{
		{
			name: "test new storage, type",
			args: args{
				storage: &repository.Maps{},
				log:     &logger.ZapLogger{},
			},
			want: &StorageService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorageService(tt.args.log, tt.args.storage); !reflect.DeepEqual(reflect.TypeOf(got).Name(), reflect.TypeOf(tt.want).Name()) {
				t.Errorf("NewStorageService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageService_SenderGetValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	want := map[string]any{"lastgc": 1.5, "count": int64(7)}

	st.EXPECT().GetValues().Return(want, nil)

	got, err := svc.SenderGetValues()
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestStorageService_SenderGetValue_Gauge_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	st.EXPECT().GetValueGauge("lastgc").Return(12.34, nil)

	res, err := svc.SenderGetValue(models.GetValueRequest{
		ID:    "lastgc",
		MType: models.Gauge,
	})
	require.NoError(t, err)
	require.Equal(t, "lastgc", res.ID)
	require.Equal(t, models.Gauge, res.MType)
	require.NotNil(t, res.Value)
	require.Equal(t, 12.34, *res.Value)
	require.Nil(t, res.Delta)
}

func TestStorageService_SenderGetValue_Counter_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	st.EXPECT().GetValueCounter("count").Return(int64(42), nil)

	res, err := svc.SenderGetValue(models.GetValueRequest{
		ID:    "count",
		MType: models.Counter,
	})
	require.NoError(t, err)
	require.Equal(t, "count", res.ID)
	require.Equal(t, models.Counter, res.MType)
	require.NotNil(t, res.Delta)
	require.Equal(t, int64(42), *res.Delta)
	require.Nil(t, res.Value)
}

func TestStorageService_SenderGetValue_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	wantErr := errors.New("down")
	st.EXPECT().GetValueGauge("lastgc").Return(0.0, wantErr)

	_, err := svc.SenderGetValue(models.GetValueRequest{
		ID:    "lastgc",
		MType: models.Gauge,
	})
	require.ErrorIs(t, err, wantErr)
}

func TestStorageService_SenderGetValue_BadType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	_, err := svc.SenderGetValue(models.GetValueRequest{
		ID:    "x",
		MType: "xx",
	})
	require.Error(t, err)
}

func TestStorageService_SenderPostUpdate_Gauge_FromValueStr_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	st.EXPECT().SaverValue("lastgc", 1.25).Return(nil)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:       "lastgc",
		MType:    models.Gauge,
		ValueStr: "1.25",
	})
	require.NoError(t, err)
}

func TestStorageService_SenderPostUpdate_Gauge_FromValuePtr_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	v := 9.99
	st.EXPECT().SaverValue("lastgc", 9.99).Return(nil)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:    "lastgc",
		MType: models.Gauge,
		Value: &v,
	})
	require.NoError(t, err)
}

func TestStorageService_SenderPostUpdate_Gauge_ParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:       "lastgc",
		MType:    models.Gauge,
		ValueStr: "notfloat",
	})
	require.Error(t, err)
}

func TestStorageService_SenderPostUpdate_Gauge_StrangeRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:    "lastgc",
		MType: models.Gauge,
	})
	require.Error(t, err)
}

func TestStorageService_SenderPostUpdate_Counter_FromValueStr_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	st.EXPECT().IncrementValue("count", int64(10)).Return(nil)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:       "count",
		MType:    models.Counter,
		ValueStr: "10",
	})
	require.NoError(t, err)
}

func TestStorageService_SenderPostUpdate_Counter_FromDeltaPtr_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	d := int64(7)
	st.EXPECT().IncrementValue("count", int64(7)).Return(nil)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:    "count",
		MType: models.Counter,
		Delta: &d,
	})
	require.NoError(t, err)
}

func TestStorageService_SenderPostUpdate_BadType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	st := mocks.NewMockMemStorage(ctrl)
	svc := NewStorageService(log, st)

	err := svc.SenderPostUpdate(models.PostUpdateRequest{
		ID:    "x",
		MType: "xx",
	})
	require.Error(t, err)
}
