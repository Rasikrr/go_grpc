package grpcapp

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
	myVal "sso/pkg/validator"
)

var val *validator.Validate

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	auth authgrpc.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	v := getValidator()

	authgrpc.Register(gRPCServer, auth, v)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("GRPC server running...", slog.Int("port", a.port))

	defer l.Close()

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping GRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}

func getValidator() *validator.Validate {
	if val != nil {
		return val
	}
	val, err := myVal.New()
	if err != nil {
		panic(err)
	}
	return val
}
