package domain

type Warehouse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	IsAvailable bool             `json:"is_available"` // забыл реализовать логику доступности/недоступности склада
	Goods       []GoodsWarehouse `json:"goods"`
}

type Good struct {
	Name       string           `json:"name"`
	Size       int              `json:"size"`
	ID         int              `json:"id"`
	Warehouses []WarehouseGoods `json:"warehouses"`
}

type WarehouseGoods struct {
	WarehouseName string `json:"name"`
	IsAvailable   bool   `json:"is_available"`
	Count         int    `json:"count"`
	Reserved      int    `json:"reserved"`
}

type GoodsWarehouse struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	ID       int    `json:"id"`
	Count    int    `json:"count"`
	Reserved int    `json:"reserved"`
}

type PairGoodWarehouse struct {
	GoodID      int   `json:"good_id"`
	WarehouseID int   `json:"warehouse_id"`
	Error       error `json:"error,omitempty"`
}

type MetaInfoReservation struct {
	ReservedPairs    []PairGoodWarehouse `json:"reserved"`
	ErrorReservation []PairGoodWarehouse `json:"error_reservation"`
}

type MetaInfoReleaseReservation struct {
	ReleasedReservations []PairGoodWarehouse `json:"released"`
	ErrorRelease         []PairGoodWarehouse `json:"error_release"`
}
