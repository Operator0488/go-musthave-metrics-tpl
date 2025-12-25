package service

import (
	"context"
	"fmt"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/repository"
	"strconv"
)

type StorageService struct {
	ctx     context.Context
	storage repository.MemStorage
}

type Service interface {
	Sender(req *models.UpdateRequest) error
}

func NewStorageService(ctx context.Context, storage repository.MemStorage) *StorageService {
	return &StorageService{
		ctx:     ctx,
		storage: storage,
	}
}

// Sender -
func (s *StorageService) Sender(req *models.UpdateRequest) error {

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
