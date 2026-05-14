package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/vpt/common/log"
	"github.com/vpt/common/middleware"
	"github.com/vpt/registry/internal/config"
	"github.com/vpt/registry/internal/handler"
	"github.com/vpt/registry/internal/repo"
	"github.com/vpt/registry/internal/service"
)

func main() {
	log.Init("registry", 0)
	cfg := config.Load()

	db, err := repo.OpenDB(cfg.DBPath)
	if err != nil {
		log.Error("open db", "err", err)
		return
	}
	defer db.Close()

	svcRepo := repo.NewServiceRepo(db)
	cfgRepo := repo.NewConfigRepo(db)

	regSvc := service.NewRegistry(svcRepo)
	cfgSvc := service.NewConfig(cfgRepo)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, regSvc, cfgSvc)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: middleware.Chain(mux, middleware.Recover, middleware.AccessLog),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go regSvc.StartHealthCheck(ctx, 10*time.Second)

	go func() {
		log.Info("registry listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "err", err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
}
