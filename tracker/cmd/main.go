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
	"github.com/vpt/tracker/internal/config"
	"github.com/vpt/tracker/internal/handler"
	"github.com/vpt/tracker/internal/repo"
	"github.com/vpt/tracker/internal/service"
	"github.com/vpt/tracker/internal/udp"
)

func main() {
	log.Init("tracker", 0)
	cfg := config.Load()

	db, err := repo.OpenDB(cfg.DBPath)
	if err != nil {
		log.Error("open db", "err", err)
		return
	}
	defer db.Close()

	peerRepo := repo.NewPeerRepo(db)
	statRepo := repo.NewStatRepo(db)

	reg := registry.NewClient(cfg.RegistryURL)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := reg.Register(ctx, registry.ServiceInfo{Name: "tracker", Address: cfg.HTTPAddr, Scheme: "http"}); err != nil {
		log.Warn("register failed", "err", err)
	}
	reg.StartHeartbeat(ctx, 10*time.Second)

	trackerSvc := service.NewTracker(peerRepo, statRepo, reg, cfg)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, trackerSvc)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: middleware.Chain(mux, middleware.Recover, middleware.AccessLog),
	}

	udpSrv := udp.NewServer(cfg.UDPAddr, trackerSvc)
	go udpSrv.Run(ctx)

	go trackerSvc.StartReaper(ctx, 60*time.Second)

	go func() {
		log.Info("tracker http listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "err", err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
}
