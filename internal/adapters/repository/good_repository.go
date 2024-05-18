package repository

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/core/domain"

	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

const getGoodById = `SELECT name, size, id FROM goods WHERE id = $1`

func (pg *PostgresConn) goodIsExist(ctx context.Context, id int) (bool, error) {
	row := pg.pool.QueryRow(ctx, getGoodById, id)

	good := domain.Good{}
	if err := row.Scan(&good.Name, &good.Size, &good.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (pg *PostgresConn) GetGood(ctx context.Context, id int) (domain.Good, error) {
	row := pg.pool.QueryRow(ctx, getGoodById, id)

	good := domain.Good{}
	if err := row.Scan(&good.Name, &good.Size, &good.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Good{}, ErrNotFound
		}
		return domain.Good{}, fmt.Errorf("error get good with id = %d: %w", id, err)
	}

	w, err := pg.getWarehousesByGoodId(ctx, id)
	if err != nil {
		return domain.Good{}, fmt.Errorf("error get warehouses with good id = %d: %w", id, err)
	}

	good.Warehouses = w

	return good, nil
}

const getWarehousesByGoodId = `SELECT warehouse.name, is_available, count, reserved FROM goods INNER JOIN goods_warehouse ON goods.id = goods_warehouse.good_id INNER JOIN warehouse ON goods_warehouse.warehouse_id = warehouse.id WHERE goods.id = $1`

func (pg *PostgresConn) getWarehousesByGoodId(ctx context.Context, id int) ([]domain.WarehouseGoods, error) {
	rows, err := pg.pool.Query(ctx, getWarehousesByGoodId, id)
	if err != nil {
		return nil, fmt.Errorf("error get warehouse by id = %d: %w", id, err)
	}

	defer rows.Close()

	ws := make([]domain.WarehouseGoods, 0)

	for rows.Next() {
		w := domain.WarehouseGoods{}

		if err = rows.Scan(&w.WarehouseName, &w.IsAvailable, &w.Count, &w.Reserved); err != nil {
			return nil, fmt.Errorf("error scan from rows: %w", err)
		}

		ws = append(ws, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return ws, nil
}

const createGood = `INSERT INTO goods VALUES ($1,$2,$3)`

func (pg *PostgresConn) CreateGood(ctx context.Context, good domain.Good) error {
	goodIsExist, err := pg.goodIsExist(ctx, good.ID)
	if err != nil {
		return fmt.Errorf("error check good is exist: %w", err)
	}
	if goodIsExist {
		return ErrIsExist
	}

	if _, err = pg.pool.Exec(ctx, createGood, good.Name, good.Size, good.ID); err != nil {
		return fmt.Errorf("error create good: %w", err)
	}

	return nil
}

const updateGood = `UPDATE goods SET name = $1, size = $2 WHERE id = $3`

func (pg *PostgresConn) UpdateGood(ctx context.Context, good domain.Good) error {
	goodIsExist, err := pg.goodIsExist(ctx, good.ID)
	if err != nil {
		return fmt.Errorf("error check good is exist: %w", err)
	}
	if !goodIsExist {
		return ErrIsNotExist
	}

	if _, err = pg.pool.Exec(ctx, updateGood, good.Name, good.Size, good.ID); err != nil {
		return fmt.Errorf("error update good: %w", err)
	}

	return nil
}

const deleteGood = `DELETE FROM goods WHERE id = $1`

func (pg *PostgresConn) DeleteGood(ctx context.Context, id int) error {
	goodIsExist, err := pg.goodIsExist(ctx, id)
	if err != nil {
		return fmt.Errorf("error check good is exist: %w", err)
	}

	if !goodIsExist {
		return ErrIsNotExist
	}

	if _, err = pg.pool.Exec(ctx, deleteGood, id); err != nil {
		return fmt.Errorf("error delete good: %w", err)
	}
	return nil
}

const checkGoodInWarehouse = `SELECT * FROM goods_warehouse WHERE warehouse_id = $1 AND good_id = $2 FOR UPDATE`

var (
	errCannotReserve = errors.New("all goods is reserved")
)

func (pg *PostgresConn) checkGoodInWarehouse(ctx context.Context, tx pgx.Tx, pair domain.PairGoodWarehouse) (bool, error) {
	row := tx.QueryRow(ctx, checkGoodInWarehouse, pair.WarehouseID, pair.GoodID)

	gw := struct {
		ID            int
		WarehouseName string
		GoodID        int
		Count         int
		Reserved      int
	}{}

	if err := row.Scan(&gw.ID, &gw.WarehouseName, &gw.GoodID, &gw.Count, &gw.Reserved); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error check good in warehouse: %w", err)
	}

	if gw.Count == gw.Reserved {
		return true, errCannotReserve
	}
	return true, nil
}

const reserve = `UPDATE goods_warehouse SET reserved = reserved + 1 WHERE warehouse_id = $1 AND good_id = $2`

func (pg *PostgresConn) Reservation(ctx context.Context, pairs []domain.PairGoodWarehouse) (domain.MetaInfoReservation, error) {
	ans := domain.MetaInfoReservation{
		ReservedPairs:    make([]domain.PairGoodWarehouse, 0, len(pairs)),
		ErrorReservation: make([]domain.PairGoodWarehouse, 0),
	}
	chReserved := AsyncWriteResult(&(ans.ReservedPairs))
	chErr := AsyncWriteResult(&(ans.ErrorReservation))
	g, gCtx := errgroup.WithContext(ctx)
	for _, pair := range pairs {
		g.Go(func() error {
			tx, err := pg.pool.BeginTx(gCtx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
			if err != nil {
				return err
			}
			defer tx.Rollback(gCtx)
			resCheck, err := pg.checkGoodInWarehouse(gCtx, tx, pair)
			if err != nil {
				if errors.Is(err, errCannotReserve) {
					pair.Error = ErrReserve
					chErr <- pair
					return nil
				}
				pair.Error = err
				chErr <- pair
				return nil
			}

			if !resCheck {
				pair.Error = ErrFailedCheckGoodInWarehouse
				chErr <- pair
				return nil
			}

			if _, err = tx.Exec(ctx, reserve, pair.WarehouseID, pair.GoodID); err != nil {
				pair.Error = err
				chErr <- pair
				return nil
			}
			if err = tx.Commit(gCtx); err != nil {
				pair.Error = err
				chErr <- pair
				return nil
			}
			chReserved <- pair
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return domain.MetaInfoReservation{}, err
	}
	close(chReserved)
	close(chErr)
	return ans, nil
}

func AsyncWriteResult(ans *[]domain.PairGoodWarehouse) chan<- domain.PairGoodWarehouse {
	ch := make(chan domain.PairGoodWarehouse)
	go func() {
		for v := range ch {
			*ans = append(*ans, v)
		}
	}()
	return ch
}

const releaseReservation = `UPDATE goods_warehouse SET reserved = 0 WHERE warehouse_id = $1 AND good_id = $2`

func (pg *PostgresConn) ReleaseReservation(ctx context.Context, pairs []domain.PairGoodWarehouse) (domain.MetaInfoReleaseReservation, error) {
	g, gCtx := errgroup.WithContext(ctx)
	ans := domain.MetaInfoReleaseReservation{
		ReleasedReservations: make([]domain.PairGoodWarehouse, 0, len(pairs)),
		ErrorRelease:         make([]domain.PairGoodWarehouse, 0),
	}
	chReleased := AsyncWriteResult(&(ans.ReleasedReservations))
	chErr := AsyncWriteResult(&(ans.ErrorRelease))
	for _, pair := range pairs {
		g.Go(func() error {
			tx, err := pg.pool.BeginTx(gCtx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
			if err != nil {
				return err
			}
			defer tx.Rollback(gCtx)
			resCheck, err := pg.checkGoodInWarehouse(ctx, tx, pair)
			if err != nil && !errors.Is(err, errCannotReserve) {
				pair.Error = err
				chErr <- pair
				return nil
			}

			if !resCheck {
				pair.Error = ErrFailedCheckGoodInWarehouse
				chErr <- pair
				return nil
			}

			if _, err = pg.pool.Exec(ctx, releaseReservation, pair.WarehouseID, pair.GoodID); err != nil {
				pair.Error = err
				chErr <- pair
				return nil
			}
			if err = tx.Commit(gCtx); err != nil {
				pair.Error = err
				chErr <- pair
				return nil
			}
			chReleased <- pair
			return nil
		})
	}
	close(chErr)
	close(chReleased)
	return ans, nil
}

const checkExistGoodOnWarehouse = `SELECT 1 FROM goods_warehouse WHERE good_id = $1 AND warehouse_id = $2`

func (pg *PostgresConn) checkIsExistGoodOnWarehouse(ctx context.Context, goodID, warehouseID int) (bool, error) {
	row := pg.pool.QueryRow(ctx, checkExistGoodOnWarehouse, goodID, warehouseID)

	if err := row.Scan(new(int)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error check exist good on warehouse: %w", err)
	}

	return true, nil
}

const addGoodOnWarehouseExist = `UPDATE goods_warehouse SET count = count + $1 WHERE warehouse_id = $2 AND good_id = $3`
const addGoodOnWarehouseNotExist = `INSERT INTO goods_warehouse(warehouse_id, good_id, count, reserved) VALUES ($1, $2, $3, 0)`

func (pg *PostgresConn) AddGoodOnWarehouse(ctx context.Context, goodID, warehouseID, count int) error {
	isExist, err := pg.warehouseIsExist(ctx, warehouseID)
	if err != nil {
		return err
	}

	if !isExist {
		return ErrIsNotExist
	}

	isExist, err = pg.goodIsExist(ctx, goodID)
	if err != nil {
		return err
	}

	if !isExist {
		return ErrIsNotExist
	}

	isExist, err = pg.checkIsExistGoodOnWarehouse(ctx, goodID, warehouseID)

	if err != nil {
		return err
	}
	if isExist {
		if _, err = pg.pool.Exec(ctx, addGoodOnWarehouseExist, count, warehouseID, goodID); err != nil {
			return err
		}
	} else {
		if _, err = pg.pool.Exec(ctx, addGoodOnWarehouseNotExist, goodID, warehouseID, count); err != nil {
			return err
		}
	}

	return nil
}
