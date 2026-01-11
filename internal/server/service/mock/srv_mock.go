package mock

import models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"

type MockService struct {
	Called bool
	Err    error
}

func (m *MockService) SenderPostUpdate(req *models.PostUpdateRequest) error {
	m.Called = true
	return m.Err
}
