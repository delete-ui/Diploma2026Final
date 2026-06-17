package service

import (
	"GolangBackendDiploma26/internal/models"
	"GolangBackendDiploma26/internal/repository"
	"context"
)

type BatteryService struct {
	repo *repository.BatteryRepository
}

func NewBatteryService(repo *repository.BatteryRepository) *BatteryService {
	return &BatteryService{repo: repo}
}

type BatteryListResponse struct {
	Batteries []models.Battery `json:"batteries"`
	Total     int              `json:"total"`
	Page      int              `json:"page"`
	Limit     int              `json:"limit"`
}

func (s *BatteryService) List(ctx context.Context, page, limit int) (*BatteryListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	batteries, total, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &BatteryListResponse{
		Batteries: batteries,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}
