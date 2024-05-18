package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"warehouse/internal/adapters/repository"
	"warehouse/internal/app"
	"warehouse/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}
	db, err := repository.NewPostgresConn(ctx, cfg.DB.String(), cfg.DB.MigrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.NewApp(db, db)
	if err != nil {
		log.Fatal(err)
	}
	if err = a.Run(ctx, cfg.Server.String()); err != nil {
		log.Fatal(err)
	}
}
