package ports

import (
	"context"
	"warehouse/internal/core/domain"
)

type GoodRepository interface {
	GetGood(ctx context.Context, id int) (domain.Good, error)
	CreateGood(ctx context.Context, good domain.Good) error
	UpdateGood(ctx context.Context, good domain.Good) error
	DeleteGood(ctx context.Context, id int) error
	Reservation(ctx context.Context, pairs []domain.PairGoodWarehouse) (domain.MetaInfoReservation, error)
	ReleaseReservation(ctx context.Context, pairs []domain.PairGoodWarehouse) (domain.MetaInfoReleaseReservation, error)
	AddGoodOnWarehouse(ctx context.Context, goodID, warehouseID, count int) error
	Close()
}

type WarehouseRepository interface {
	GetWarehouse(ctx context.Context, id int) (domain.Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse domain.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse domain.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
	GetCountGoods(ctx context.Context, id int) (int, error)
	Close()
}
