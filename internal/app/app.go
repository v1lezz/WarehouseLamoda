package app

import (
	"context"
	"net"
	"net/http"
	"time"
	"warehouse/internal/core/ports"
	"warehouse/internal/core/server"

	"golang.org/x/sync/errgroup"
)

type App struct {
	srv           *http.Server
	goodRepo      ports.GoodRepository
	warehouseRepo ports.WarehouseRepository
}

func NewApp(goodRepo ports.GoodRepository, warehouseRepo ports.WarehouseRepository) (*App, error) {
	return &App{
		goodRepo:      goodRepo,
		warehouseRepo: warehouseRepo,
	}, nil
}

func (a *App) Run(ctx context.Context, srvAddr string) error {
	g, gCtx := errgroup.WithContext(ctx)
	a.srv = server.NewServer(gCtx, a.goodRepo, a.warehouseRepo, srvAddr)
	g.Go(func() error {
		a.srv.BaseContext = func(_ net.Listener) context.Context {
			return gCtx
		}
		return a.srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return a.Close()
	})
	return g.Wait()
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}
	a.goodRepo.Close()
	a.warehouseRepo.Close()
	return nil
}
