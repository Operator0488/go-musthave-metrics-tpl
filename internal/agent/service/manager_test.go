package service

import (
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStatsManager_InitialWriteStatsFillsSomeValues(t *testing.T) {
	sm := NewStatsManager()
	require.NotNil(t, sm)

	got := sm.GetMap()

	require.Contains(t, got, model.PollCount)
	require.NotNil(t, got[model.PollCount])
	require.Equal(t, "counter", got[model.PollCount].Type)
	require.Equal(t, float64(1), got[model.PollCount].Value)

	require.Contains(t, got, model.RandomValue)
	require.NotNil(t, got[model.RandomValue])
	require.Equal(t, "gauge", got[model.RandomValue].Type)
	require.GreaterOrEqual(t, got[model.RandomValue].Value, float64(0))
	require.Less(t, got[model.RandomValue].Value, float64(1))

	require.Contains(t, got, model.Alloc)
	require.NotNil(t, got[model.Alloc])
	require.GreaterOrEqual(t, got[model.Alloc].Value, float64(0))

	require.Contains(t, got, model.Sys)
	require.NotNil(t, got[model.Sys])
	require.GreaterOrEqual(t, got[model.Sys].Value, float64(0))
}

func TestStatsManager_GetMap_ReturnsDeepCopy(t *testing.T) {
	sm := NewStatsManager()

	m1 := sm.GetMap()

	m1[model.Alloc].Value = 123456789
	m1["newkey"] = &model.Stat{Type: "gauge", Value: 999}

	m2 := sm.GetMap()

	require.NotContains(t, m2, "newkey")
	require.NotNil(t, m2[model.Alloc])
	require.NotEqual(t, float64(123456789), m2[model.Alloc].Value)

	require.NotSame(t, m1[model.Alloc], m2[model.Alloc])
}
