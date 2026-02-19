package repository

import (
	"context"
	mocklog "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger/mocks"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile_CreatesDirsAndFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "dir1", "dir2", "file.txt")

	f, err := loadFile(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	defer f.Close()

	_, statErr := os.Stat(path)
	require.NoError(t, statErr)
}

func TestMaps_SaverValue_NewGauge_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	err = m.SaveValue(context.Background(), "LastGc", 12.34)
	require.NoError(t, err)

	got, err := m.GetValueGauge(context.Background(), "LastGc")
	require.NoError(t, err)
	require.Equal(t, 12.34, got)
}

func TestMaps_SaverValue_TypeMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	// положим counter под тем же именем
	m.storage["LastGc"] = models.MetricStore{
		ID:    "LastGc",
		MType: models.Counter,
		Delta: 10,
	}

	err = m.SaveValue(context.Background(), "LastGc", 1.0)
	require.ErrorIs(t, err, models.ErrorDiffType)
}

func TestMaps_IncrementValue_NewCounter_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	require.NoError(t, m.IncrementValue(context.Background(), "Lastgc", 10))
	require.NoError(t, m.IncrementValue(context.Background(), "Lastgc", 5))

	got, err := m.GetValueCounter(context.Background(), "Lastgc")
	require.NoError(t, err)
	require.Equal(t, int64(15), got)
}

func TestMaps_IncrementValue_TypeMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	m.storage["Lastgc"] = models.MetricStore{
		ID:    "Lastgc",
		MType: models.Gauge,
		Value: 1.23,
	}

	err = m.IncrementValue(context.Background(), "Lastgc", 1)
	require.ErrorIs(t, err, models.ErrorGetValue)
}

func TestMaps_GetValue_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	_, err = m.GetValueGauge(context.Background(), "Lastgc")
	require.ErrorIs(t, err, models.ErrorNotDB)

	_, err = m.GetValueCounter(context.Background(), "Lastgc")
	require.ErrorIs(t, err, models.ErrorNotDB)
}

func TestMaps_GetValues_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	require.NoError(t, m.SaveValue(context.Background(), "Lastgc", 1.5))
	require.NoError(t, m.IncrementValue(context.Background(), "Lastgc2", 7))

	all, err := m.GetValues(context.Background())
	require.NoError(t, err)

	require.Equal(t, 1.5, all["Lastgc"])
	require.Equal(t, int64(7), all["Lastgc2"])
}

func TestMaps_GetValues_UnknownType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	defer m.file.Close()

	m.storage["Lastgc"] = models.MetricStore{
		ID:    "Lastgc",
		MType: "Lastgc",
	}

	_, err = m.GetValues(context.Background())
	require.ErrorIs(t, err, models.ErrorUnType)
}

func TestMaps_PersistAndRestore_FromFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "file.txt")

	log := mocklog.NewMockLogger(ctrl)

	log.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	m1, err := NewMaps(context.Background(), log, false, path, 0)
	require.NoError(t, err)
	require.NoError(t, m1.SaveValue(context.Background(), "Lastgc", 9.99))
	require.NoError(t, m1.SaveValue(context.Background(), "Lastgc", 5.99))
	require.NoError(t, m1.IncrementValue(context.Background(), "Lastgc2", 10))
	require.NoError(t, m1.IncrementValue(context.Background(), "Lastgc2", 5))
	require.NoError(t, m1.file.Close())

	gotCounter, err := m1.GetValueCounter(context.Background(), "Lastgc2")
	require.NoError(t, err)
	require.Equal(t, int64(15), gotCounter)

	gotGauge, err := m1.GetValueGauge(context.Background(), "Lastgc")
	require.NoError(t, err)
	require.Equal(t, 5.99, gotGauge)

	m2, err := NewMaps(context.Background(), log, true, path, 0)
	require.NoError(t, err)
	defer m2.file.Close()

	gotGauge, err = m2.GetValueGauge(context.Background(), "Lastgc")
	require.NoError(t, err)
	require.Equal(t, 5.99, gotGauge)

	gotCounter, err = m2.GetValueCounter(context.Background(), "Lastgc2")
	require.NoError(t, err)
	require.Equal(t, int64(15), gotCounter)
}
