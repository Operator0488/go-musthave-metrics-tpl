package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/logger"
	models "github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/model"
	"github.com/Operator0488/go-musthave-metrics-tpl.git/internal/server/repository"
)

type StorageService struct {
	storage repository.MemStorage
	log     logger.Logger
}

//go:generate mockgen -source=srv.go -destination=./mocks/mock_srv.go -package=mocks
type Service interface {
	SenderPostUpdates(ctx context.Context, req []models.PostUpdateRequest) error
	SenderPostUpdate(ctx context.Context, req models.PostUpdateRequest) error
	SenderGetValue(ctx context.Context, req models.GetValueRequest) (models.GetValueResponse, error)
	SenderGetValues(ctx context.Context) (map[string]any, error)
	PingDB(ctx context.Context) error
}

func NewStorageService(log logger.Logger, storage repository.MemStorage) *StorageService {
	return &StorageService{
		storage: storage,
		log:     log,
	}
}

func (s *StorageService) SenderGetValues(ctx context.Context) (map[string]any, error) {
	return s.storage.GetValues(ctx)
}

func (s *StorageService) SenderGetValue(ctx context.Context, req models.GetValueRequest) (models.GetValueResponse, error) {
	res := models.GetValueResponse{
		ID:    req.ID,
		MType: req.MType,
	}

	switch req.MType {

	case models.Gauge:
		num, err := s.storage.GetValueGauge(ctx, req.ID)
		if err != nil {
			return res, err
		}
		res.Value = &num
		return res, nil

	case models.Counter:
		num, err := s.storage.GetValueCounter(ctx, req.ID)
		if err != nil {
			return res, err
		}
		res.Delta = &num
		return res, nil

	default:
		return res, fmt.Errorf("ошибка, нет подходящего типа")
	}
}

func (s *StorageService) SenderPostUpdate(ctx context.Context, req models.PostUpdateRequest) error {

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
			err := s.storage.SaveValue(ctx, req.ID, *req.Value)
			return err
		}
		return fmt.Errorf("ошибка, странный запрос")

	case models.Counter:
		if req.ValueStr != "" && req.Delta == nil {
			val, err := getValueInt(req.ValueStr)
			if err != nil {
				return err
			}
			req.Delta = &val
		}
		if req.Delta != nil {
			err := s.storage.IncrementValue(ctx, req.ID, *req.Delta)
			return err
		}
		return fmt.Errorf("ошибка, странный запрос")

	default:
		return fmt.Errorf("ошибка, нет подходящего типа")

	}
}

func (s *StorageService) SenderPostUpdates(ctx context.Context, reqs []models.PostUpdateRequest) error {
	var errs []error

	valid := make([]models.PostUpdateRequest, 0, len(reqs))

	for i, req := range reqs {
		var err error

		switch req.MType {

		case models.Gauge:
			if req.ValueStr != "" && req.Value == nil {
				val, e := getValueFloat(req.ValueStr)
				if e != nil {
					err = e
					reqs = append(reqs[:i], reqs[i+1:]...)
					continue
				}
				req.Value = &val
			}
			valid = append(valid, req)

		case models.Counter:
			if req.ValueStr != "" && req.Delta == nil {
				val, e := getValueInt(req.ValueStr)
				if e != nil {
					err = e
					continue
				}
				req.Delta = &val
			}
			valid = append(valid, req)

		default:
			err = fmt.Errorf("ошибка, нет подходящего типа")

		}

		if err != nil {
			errs = append(errs, fmt.Errorf("не записаны %q (%s) ошибка: %w", req.ID, req.MType, err))
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	err := s.storage.SaveValues(ctx, valid)
	if err != nil {
		return fmt.Errorf("не удалось сохранить данные: %w", err)
	}
	return nil
}

func (s *StorageService) PingDB(ctx context.Context) error {
	return s.storage.PingDB(ctx)
}
