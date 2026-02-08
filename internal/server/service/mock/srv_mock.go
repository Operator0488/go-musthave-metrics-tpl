package mock

import (
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"reflect"
)

type MockService struct {
	Called bool
	Err    error
	Res    any
}

func (m *MockService) SenderPostUpdate(req models.PostUpdateRequest) error {
	m.Called = true
	return m.Err
}

func (m *MockService) SenderGetValue(req models.GetValueRequest) (models.GetValueResponse, error) {
	m.Called = true
	v := reflect.ValueOf(m.Res)
	rawVal := v.Interface()
	res, ok := rawVal.(models.GetValueResponse)
	if ok {
		return models.GetValueResponse{}, fmt.Errorf("Что-то не так с передаваемым ответом")
	}
	return res, m.Err
}

func (m *MockService) SenderGetValues() (map[string]any, error) {
	m.Called = true
	return make(map[string]interface{}), m.Err
}
