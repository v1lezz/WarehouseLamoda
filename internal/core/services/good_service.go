package services

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/adapters/repository"
	"warehouse/internal/core/domain"
	"warehouse/internal/core/ports"
)

type GoodService struct {
	repo ports.GoodRepository
}

func NewGoodService(repo ports.GoodRepository) *GoodService {
	return &GoodService{repo: repo}
}

func (gs *GoodService) validateID(id int) bool {
	return id > 0
}

func (gs *GoodService) GetGood(ctx context.Context, id int) (domain.Good, error) {
	if !gs.validateID(id) {
		return domain.Good{}, ErrGoodIDisNegative
	}
	good, err := gs.repo.GetGood(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Good{}, ErrGoodNotFound
		}
		return domain.Good{}, fmt.Errorf("error get good: %w", err)
	}
	return good, nil
}

func (gs *GoodService) validateGood(good domain.Good) bool {
	return good.ID > 0 && good.Name != "" && good.Size != 0
}

func (gs *GoodService) CreateGood(ctx context.Context, good domain.Good) error {
	if !gs.validateGood(good) {
		return ErrInvalidGood
	}
	if err := gs.repo.CreateGood(ctx, good); err != nil {
		if errors.Is(err, repository.ErrIsExist) {
			return ErrGoodIsExist
		}
		return fmt.Errorf("error create good: %w", err)
	}
	return nil
}

func (gs *GoodService) UpdateGood(ctx context.Context, good domain.Good) error {
	if !gs.validateGood(good) {
		return ErrInvalidGood
	}

	if err := gs.repo.UpdateGood(ctx, good); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrGoodIsNotExist
		}
		return fmt.Errorf("error update good: %w", err)
	}
	return nil
}

func (gs *GoodService) DeleteGood(ctx context.Context, id int) error {
	if !gs.validateID(id) {
		return ErrGoodIDisNegative
	}

	if err := gs.repo.DeleteGood(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrGoodIsNotExist
		}
		return fmt.Errorf("error delete good: %w", err)
	}

	return nil
}

func (gs *GoodService) filterPairs(pairs []domain.PairGoodWarehouse) ([]domain.PairGoodWarehouse, []domain.PairGoodWarehouse) {
	filteredPairs := make([]domain.PairGoodWarehouse, 0, len(pairs))
	errPairs := make([]domain.PairGoodWarehouse, 0)
	for _, pair := range pairs {
		if !gs.validateID(pair.GoodID) {
			pair.Error = ErrGoodIDisNegative
			errPairs = append(errPairs, pair)
			continue
		}

		if !gs.validateID(pair.WarehouseID) {
			pair.Error = ErrWarehouseIDisNegative
			errPairs = append(errPairs, pair)
			continue
		}

		filteredPairs = append(filteredPairs, pair)
	}
	return filteredPairs, errPairs
}

func (gs *GoodService) Reserve(ctx context.Context, pairs []domain.PairGoodWarehouse) (domain.MetaInfoReservation, error) {
	filteredPairs, errPairs := gs.filterPairs(pairs)
	res, err := gs.repo.Reservation(ctx, filteredPairs)
	if err != nil {
		return domain.MetaInfoReservation{}, fmt.Errorf("error reserve: %w", err)
	}
	res.ErrorReservation = append(res.ErrorReservation, errPairs...)
	return res, nil
}

func (gs *GoodService) ReleaseReservation(ctx context.Context, pairs []domain.PairGoodWarehouse) (domain.MetaInfoReleaseReservation, error) {
	filteredPairs, errPairs := gs.filterPairs(pairs)
	res, err := gs.repo.ReleaseReservation(ctx, filteredPairs)
	if err != nil {
		return domain.MetaInfoReleaseReservation{}, fmt.Errorf("error release reserve: %w", err)
	}
	res.ErrorRelease = append(res.ErrorRelease, errPairs...)
	return res, nil
}

func (gs *GoodService) AddGoodOnWarehouse(ctx context.Context, goodID, warehouseID, count int) error {
	if !gs.validateID(goodID) {
		return ErrGoodIDisNegative
	}

	if !gs.validateID(warehouseID) {
		return ErrWarehouseIDisNegative
	}

	if !gs.validateID(count) {
		return ErrCountIsNegative
	}

	if err := gs.repo.AddGoodOnWarehouse(ctx, goodID, warehouseID, count); err != nil {
		return err
	}
	return nil
}
