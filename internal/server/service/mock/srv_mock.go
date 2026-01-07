package mock

import models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"

type MockService struct {
	Called bool
	Err    error
}

func (m *MockService) Sender(req *models.UpdateRequest) error {
	m.Called = true
	return m.Err
}
