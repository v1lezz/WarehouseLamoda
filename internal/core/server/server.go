package server

import (
	"context"
	"net"
	"net/http"
	"warehouse/internal/adapters/handler"
	"warehouse/internal/core/ports"
	"warehouse/internal/core/services"
)

func NewServer(ctx context.Context, goodRepo ports.GoodRepository, warehouseRepo ports.WarehouseRepository, srvAddr string) *http.Server {
	router := http.NewServeMux()

	goodService := services.NewGoodService(goodRepo)
	goodHandler := handler.NewGoodHandler(*goodService)

	router.HandleFunc("GET /getGood", goodHandler.GetGood)
	router.HandleFunc("POST /createGood", goodHandler.CreateGood)
	router.HandleFunc("PUT /updateGood", goodHandler.UpdateGood)
	router.HandleFunc("DELETE /deleteGood", goodHandler.DeleteGood)
	router.HandleFunc("PATCH /reserveGood", goodHandler.ReserveGood)
	router.HandleFunc("PATCH /releaseReservationGood", goodHandler.ReleaseReservationGood)
	router.HandleFunc("POST /addGoodOnWarehouse", goodHandler.AddGoodOnWarehouse)

	warehouseService := services.NewWarehouseService(warehouseRepo)
	warehouseHandler := handler.NewWarehouseHandler(*warehouseService)

	router.HandleFunc("GET /getWarehouse", warehouseHandler.GetWarehouse)
	router.HandleFunc("POST /createWarehouse", warehouseHandler.CreateWarehouse)
	router.HandleFunc("PUT /updateWarehouse", warehouseHandler.UpdateWarehouse)
	router.HandleFunc("DELETE /deleteWarehouse", warehouseHandler.DeleteWarehouse)
	router.HandleFunc("GET /getCountGoods", warehouseHandler.GetCountGoods)

	return &http.Server{
		Addr:    srvAddr,
		Handler: router,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}
