package services

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/adapters/repository"
	"warehouse/internal/core/domain"
	"warehouse/internal/core/ports"
)

type WarehouseService struct {
	repo ports.WarehouseRepository
}

func NewWarehouseService(repo ports.WarehouseRepository) *WarehouseService {
	return &WarehouseService{
		repo: repo,
	}
}

func (ws *WarehouseService) validateID(ID int) bool {
	return ID > 0
}

func (ws *WarehouseService) validateWarehouse(warehouse domain.Warehouse) bool {
	return warehouse.ID > 0 && warehouse.Name != ""
}

func (ws *WarehouseService) GetWarehouse(ctx context.Context, warehouseID int) (domain.Warehouse, error) {
	if !ws.validateID(warehouseID) {
		return domain.Warehouse{}, ErrWarehouseIDisNegative
	}

	warehouse, err := ws.repo.GetWarehouse(ctx, warehouseID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Warehouse{}, ErrWarehouseNotFound
		}
		return domain.Warehouse{}, fmt.Errorf("erorr get warehouse: %w", err)
	}

	return warehouse, nil
}

func (ws *WarehouseService) CreateWarehouse(ctx context.Context, warehouse domain.Warehouse) error {
	if !ws.validateWarehouse(warehouse) {
		return ErrInvalidWarehouse
	}

	if err := ws.repo.CreateWarehouse(ctx, warehouse); err != nil {
		if errors.Is(err, repository.ErrIsExist) {
			return ErrWarehouseIsExist
		}
		return fmt.Errorf("error create warehouse: %w", err)
	}

	return nil
}

func (ws *WarehouseService) UpdateWarehouse(ctx context.Context, warehouse domain.Warehouse) error {
	if !ws.validateWarehouse(warehouse) {
		return ErrInvalidWarehouse
	}

	if err := ws.repo.UpdateWarehouse(ctx, warehouse); err != nil {
		if errors.Is(err, repository.ErrIsNotExist) {
			return ErrWarehouseIsNotExist
		}
	}

	return nil
}

func (ws *WarehouseService) DeleteWarehouse(ctx context.Context, warehouseID int) error {
	if !ws.validateID(warehouseID) {
		return ErrWarehouseIDisNegative
	}

	if err := ws.repo.DeleteWarehouse(ctx, warehouseID); err != nil {
		if errors.Is(err, repository.ErrIsNotExist) {
			return ErrWarehouseIsNotExist
		}
	}

	return nil
}

func (ws *WarehouseService) GetCountGoodsByWarehouseID(ctx context.Context, warehouseID int) (int, error) {
	if !ws.validateID(warehouseID) {
		return 0, ErrWarehouseIDisNegative
	}

	cnt, err := ws.repo.GetCountGoods(ctx, warehouseID)
	if err != nil {
		if errors.Is(err, repository.ErrIsNotExist) {
			return 0, ErrWarehouseIsNotExist
		}

		return 0, fmt.Errorf("error get count goods on warehouse: %w", err)
	}
	return cnt, nil
}
