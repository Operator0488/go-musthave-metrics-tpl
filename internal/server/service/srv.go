package service

import (
	"context"
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
	"log"
	"strconv"
)

type StorageService struct {
	ctx     context.Context
	storage repository.MemStorage
}

type Service interface {
	SenderPostUpdate(req *models.PostUpdateRequest) error
	SenderGetValue(req *models.GetValueRequest) (string, error)
	SenderGetValues() (map[string]interface{}, error)
}

func NewStorageService(ctx context.Context, storage repository.MemStorage) *StorageService {
	return &StorageService{
		ctx:     ctx,
		storage: storage,
	}
}

// SenderGetValues -
func (s *StorageService) SenderGetValues() (map[string]interface{}, error) {
	return s.storage.GetValues()
}

// SenderGetValue -
func (s *StorageService) SenderGetValue(req *models.GetValueRequest) (string, error) {
	switch req.Type {

	case models.Gauge:
		log.Println(req.Name)
		num, err := s.storage.GetValueGauge(req.Name)
		if err != nil {
			return "", err
		}
		str := getStringFloat(num)
		if str == "" {
			return "", fmt.Errorf("Ошибка получения значения")
		}
		return str, nil

	case models.Counter:
		num, err := s.storage.GetValueCounter(req.Name)
		if err != nil {
			return "", err
		}
		str := getStringInt(num)
		if str == "" {
			return "", fmt.Errorf("Ошибка получения значения")
		}
		return str, nil

	default:
		return "", fmt.Errorf("Ошибка, нет подходящего типа")

	}
}

// SenderPostUpdate -
func (s *StorageService) SenderPostUpdate(req *models.PostUpdateRequest) error {

	switch req.Type {

	case models.Gauge:
		val, err := getValueFloat(req.Value)
		if err != nil {
			return err
		}
		err = s.storage.SaverValue(req.Name, val)
		return err

	case models.Counter:
		val, err := getValueInt(req.Value)
		if err != nil {
			return err
		}
		err = s.storage.IncrementValue(req.Name, val)
		return err

	default:
		return fmt.Errorf("Ошибка, нет подходящего типа")

	}

}

func getValueFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func getValueInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func getStringFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func getStringInt(i int64) string {
	return strconv.FormatInt(i, 10)
}
