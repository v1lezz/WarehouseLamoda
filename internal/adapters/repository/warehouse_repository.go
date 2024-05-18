package repository

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/core/domain"

	"github.com/jackc/pgx/v5"
)

const getWarehouse = `SELECT * FROM warehouse WHERE id = $1`

func (pg *PostgresConn) GetWarehouse(ctx context.Context, id int) (domain.Warehouse, error) {
	row := pg.pool.QueryRow(ctx, getWarehouse, id)

	w := domain.Warehouse{}
	if err := row.Scan(&w.ID, &w.Name, &w.IsAvailable); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Warehouse{}, ErrNotFound
		}
		return domain.Warehouse{}, fmt.Errorf("error get warehouse with id = %d: %w", id, err)
	}

	gs, err := pg.getGoodsByWarehouseId(ctx, id)
	if err != nil {
		return domain.Warehouse{}, fmt.Errorf("error get goods by warehouse id = %d: %w", id, err)
	}

	w.Goods = gs

	return w, nil
}

const getGoodsByWarehouseId = `SELECT goods.name, size, goods.id, count, reserved FROM goods INNER JOIN goods_warehouse ON goods.id = goods_warehouse.good_id INNER JOIN warehouse ON goods_warehouse.warehouse_id = warehouse.id WHERE goods.id = $1`

func (pg *PostgresConn) getGoodsByWarehouseId(ctx context.Context, id int) ([]domain.GoodsWarehouse, error) {
	rows, err := pg.pool.Query(ctx, getGoodsByWarehouseId, id)
	if err != nil {
		return nil, fmt.Errorf("error get warehouse by id = %d: %w", id, err)
	}

	defer rows.Close()

	gs := make([]domain.GoodsWarehouse, 0)

	for rows.Next() {
		g := domain.GoodsWarehouse{}

		if err = rows.Scan(&g.Name, &g.Size, &g.ID, &g.Count, &g.Reserved); err != nil {
			return nil, fmt.Errorf("error scan from rows: %w", err)
		}

		gs = append(gs, g)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows errors: %w", err)
	}

	return gs, nil
}

const checkWarehouse = `SELECT 1 FROM warehouse WHERE id = $1`

func (pg *PostgresConn) warehouseIsExist(ctx context.Context, id int) (bool, error) {
	row := pg.pool.QueryRow(ctx, checkWarehouse, id)

	if err := row.Scan(new(int)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error check exist warehouse with id = %d: %w", id, err)
	}

	return true, nil
}

const createWarehouse = `INSERT INTO warehouse(name, is_available) VALUES ($1, $2)`

func (pg *PostgresConn) CreateWarehouse(ctx context.Context, warehouse domain.Warehouse) error {
	isExist, err := pg.warehouseIsExist(ctx, warehouse.ID)
	if err != nil {
		return fmt.Errorf("error check warehouse is exist: %w", err)
	}

	if isExist {
		return ErrIsExist
	}

	if _, err = pg.pool.Exec(ctx, createWarehouse, warehouse.Name, warehouse.IsAvailable); err != nil {
		return fmt.Errorf("error create warehouse: %w", err)
	}
	return nil
}

const updateWarehouse = `UPDATE warehouse SET name = $1, is_available = $2 WHERE id = $3`

func (pg *PostgresConn) UpdateWarehouse(ctx context.Context, warehouse domain.Warehouse) error {
	isExist, err := pg.warehouseIsExist(ctx, warehouse.ID)
	if err != nil {
		return fmt.Errorf("error check warehouse is exist: %w", err)
	}

	if !isExist {
		return ErrIsExist
	}

	if _, err = pg.pool.Exec(ctx, updateWarehouse, warehouse.Name, warehouse.IsAvailable, warehouse.ID); err != nil {
		return fmt.Errorf("error update warehouse: %w")
	}

	return nil
}

const deleteWarehouse = `DELETE FROM warehouse WHERE id = $1`

func (pg *PostgresConn) DeleteWarehouse(ctx context.Context, id int) error {
	isExist, err := pg.warehouseIsExist(ctx, id)
	if err != nil {
		return fmt.Errorf("error check warehouse is exist: %w", err)
	}

	if !isExist {
		return ErrIsExist
	}

	if _, err = pg.pool.Exec(ctx, deleteWarehouse, id); err != nil {
		return fmt.Errorf("error delete warehouse: %w", err)
	}

	return nil
}

const getCountGoods = `SELECT SUM(count) FROM goods INNER JOIN goods_warehouse ON goods.id = goods_warehouse.good_id INNER JOIN warehouse ON goods_warehouse.warehouse_id = warehouse.id WHERE warehouse.id = $1 GROUP BY warehouse.id`

func (pg *PostgresConn) GetCountGoods(ctx context.Context, id int) (int, error) {
	isExist, err := pg.warehouseIsExist(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("error check warehouse is exist: %w", err)
	}

	if !isExist {
		return 0, ErrIsNotExist
	}

	row := pg.pool.QueryRow(ctx, getCountGoods, id)
	cnt := 0
	if err = row.Scan(&cnt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("error get count goods in warehouse with id = %d: %w", id, err)
	}

	return cnt, nil
}
