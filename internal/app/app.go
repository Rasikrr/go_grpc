package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	auth2 "sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

// TODO: Storage

func New(
	log *slog.Logger,
	storageCfg *config.StorageConfig,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	// Create storage
	storage, err := postgres.New(storageCfg)
	if err != nil {
		panic(err)
	}
	auth := auth2.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, auth, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
