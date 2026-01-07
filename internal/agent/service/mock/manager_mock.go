package mock

import (
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/agent/model"
)

type MockManager struct {
	M map[string]*model.Stat
}

func (mm *MockManager) GetMap() map[string]*model.Stat {
	return mm.M
}

func (mm *MockManager) WriteStats() {

}
