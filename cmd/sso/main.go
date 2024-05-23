package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Run app: go run cmd/sso/main.go --config=./cmd/config/local.yaml

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting application",
		slog.String("env", cfg.Env),
	)

	application := app.New(log, &cfg.Storage, cfg.GRPC.Port, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	application.GRPCSrv.Stop()
	log.Info("application stopped", slog.Any("signal", sign))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					AddSource: true,
					Level:     slog.LevelDebug,
				}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{
					AddSource: true,
					Level:     slog.LevelDebug,
				}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					AddSource: true,
					Level:     slog.LevelInfo,
				}),
		)
	}
	return log
}
