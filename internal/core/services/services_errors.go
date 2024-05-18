package services

import "errors"

var (
	ErrGoodNotFound            = errors.New("good with this id is not found")
	ErrGoodIDisNegative        = errors.New("good id is negative")
	ErrWarehouseIDisNegative   = errors.New("warehouse id is negative")
	ErrGoodIsExist             = errors.New("good with this id already exist")
	ErrInvalidGood             = errors.New("good is invalid")
	ErrGoodIsNotExist          = errors.New("good with this id is not exist")
	ErrGoodWarehouseIsNotExist = errors.New("good or warehouse with this id is not exist")
	ErrReservation             = errors.New("good in this warehouse is not exist")
	ErrWarehouseNotFound       = errors.New("warehouse with this id is not found")
	ErrInvalidWarehouse        = errors.New("warehouse is invalid")
	ErrWarehouseIsExist        = errors.New("warehouse whit this id is exist")
	ErrWarehouseIsNotExist     = errors.New("warehouse with this id is not exist")
	ErrCountIsNegative         = errors.New("count is negative")
)
