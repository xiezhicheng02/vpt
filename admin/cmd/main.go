package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/vpt/common/log"
	"github.com/vpt/common/middleware"
	"github.com/vpt/common/registry"
	"github.com/vpt/admin/internal/config"
	"github.com/vpt/admin/internal/handler"
	"github.com/vpt/admin/internal/repo"
	"github.com/vpt/admin/internal/service"
)

func main() {
	log.Init("admin", 0)
	cfg := config.Load()

	db, err := repo.OpenDB(cfg.DBPath)
	if err != nil {
		log.Error("open db", "err", err)
		return
	}
	defer db.Close()

	userRepo := repo.NewUserRepo(db)
	torrentRepo := repo.NewTorrentRepo(db)
	tokenRepo := repo.NewTokenRepo(db)

	userSvc := service.NewUser(userRepo, tokenRepo, cfg)
	torrentSvc := service.NewTorrent(torrentRepo, cfg)

	reg := registry.NewClient(cfg.RegistryURL)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := reg.Register(ctx, registry.ServiceInfo{Name: "admin", Address: cfg.Addr, Scheme: "http"}); err != nil {
		log.Warn("register failed", "err", err)
	}
	reg.StartHeartbeat(ctx, 10*time.Second)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, userSvc, torrentSvc)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: middleware.Chain(mux, middleware.Recover, middleware.AccessLog),
	}

	go func() {
		log.Info("admin listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "err", err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
}
