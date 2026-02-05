package service

import (
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
)

type StorageService struct {
	storage repository.MemStorage
}

type Service interface {
	SenderPostUpdate(req models.PostUpdateRequest) error
	SenderGetValue(req models.GetValueRequest) (models.GetValueResponse, error)
	SenderGetValues() (map[string]any, error)
}

func NewStorageService(storage repository.MemStorage) *StorageService {
	return &StorageService{
		storage: storage,
	}
}

func (s *StorageService) SenderGetValues() (map[string]any, error) {
	return s.storage.GetValues()
}

func (s *StorageService) SenderGetValue(req models.GetValueRequest) (models.GetValueResponse, error) {
	res := models.GetValueResponse{
		ID:    req.ID,
		MType: req.MType,
	}

	switch req.MType {

	case models.Gauge:
		num, err := s.storage.GetValueGauge(req.ID)
		if err != nil {
			return res, err
		}
		res.Value = &num
		return res, nil

	case models.Counter:
		num, err := s.storage.GetValueCounter(req.ID)
		if err != nil {
			return res, err
		}
		res.Delta = &num
		return res, nil

	default:
		return res, fmt.Errorf("Ошибка, нет подходящего типа")
	}
}

func (s *StorageService) SenderPostUpdate(req models.PostUpdateRequest) error {

	switch req.MType {

	case models.Gauge:
		if req.ValueStr != "" && req.Value == nil {
			val, err := getValueFloat(req.ValueStr)
			if err != nil {
				return err
			}
			req.Value = &val
		}
		if req.Value != nil {
			err := s.storage.SaverValue(req.ID, *req.Value)
			return err
		}
		return fmt.Errorf("Ошибка, странный запрос")

	case models.Counter:
		if req.ValueStr != "" && req.Delta == nil {
			val, err := getValueInt(req.ValueStr)
			if err != nil {
				return err
			}
			req.Delta = &val
		}
		if req.Delta != nil {
			err := s.storage.IncrementValue(req.ID, *req.Delta)
			return err
		}
		return fmt.Errorf("Ошибка, странный запрос")

	default:
		return fmt.Errorf("Ошибка, нет подходящего типа")

	}
}
