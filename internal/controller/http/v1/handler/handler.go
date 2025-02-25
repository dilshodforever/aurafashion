package handler

import (
	"aura-fashion/config"
	"aura-fashion/internal/usecase"
	"aura-fashion/pkg/logger"

	rediscache "github.com/golanguzb70/redis-cache"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	MinIO   *minio.Client
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	Redis   rediscache.RedisCache
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache,MinIO *minio.Client) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		Redis:   redis,
		MinIO:     MinIO,
	}
}
