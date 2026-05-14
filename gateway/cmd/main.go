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
	"github.com/vpt/gateway/internal/auth"
	"github.com/vpt/gateway/internal/config"
	"github.com/vpt/gateway/internal/handler"
	"github.com/vpt/gateway/internal/proxy"
	"github.com/vpt/gateway/internal/udp"
)

func main() {
	log.Init("gateway", 0)
	cfg := config.Load()

	reg := registry.NewClient(cfg.RegistryURL)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := reg.Register(ctx, registry.ServiceInfo{Name: "gateway", Address: cfg.HTTPAddr, Scheme: "http"}); err != nil {
		log.Warn("register failed", "err", err)
	}
	reg.StartHeartbeat(ctx, 10*time.Second)

	authClient := auth.NewClient(reg)
	router := proxy.NewRouter(reg)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, router, authClient)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: middleware.Chain(mux, middleware.Recover, middleware.AccessLog),
	}

	udpSrv := udp.NewServer(cfg.UDPAddr, router)
	go udpSrv.Run(ctx)

	go func() {
		log.Info("gateway http listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "err", err)
		}
	}()

	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
}
