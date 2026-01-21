package mock

import models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"

type MockService struct {
	Called bool
	Err    error
	Str    string
}

func (m *MockService) SenderPostUpdate(req *models.PostUpdateRequest) error {
	m.Called = true
	return m.Err
}

func (m *MockService) SenderGetValue(req *models.GetValueRequest) (string, error) {
	m.Called = true
	return m.Str, m.Err
}

func (m *MockService) SenderGetValues() (map[string]interface{}, error) {
	m.Called = true
	return make(map[string]interface{}), m.Err
}
