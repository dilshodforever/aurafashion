// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"aura-fashion/config"
	v1 "aura-fashion/internal/controller/http/v1"
	"aura-fashion/internal/usecase"
	"aura-fashion/pkg/httpserver"
	"aura-fashion/pkg/logger"
	"aura-fashion/pkg/postgres"

	rediscache "github.com/golanguzb70/redis-cache"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	
	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	useCase := usecase.New(pg, cfg, l)

	// redis
	redis, err := rediscache.New(&rediscache.Config{
		RedisHost: cfg.Redis.RedisHost,
		RedisPort: cfg.Redis.RedisPort,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rediscache.New: %w", err))
	}
	
	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, cfg, useCase, redis)
	
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	
	l.Info(fmt.Sprintf("app - Run - httpServer: %s", cfg.HTTP.Port))
	
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	
}
