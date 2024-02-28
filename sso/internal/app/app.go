package app

import (
	"time"

	grpcapp "github.com/Faner201/go-grpc-server/sso/internal/app/grpc"
	"github.com/Faner201/go-grpc-server/sso/internal/config"
	"github.com/Faner201/go-grpc-server/sso/internal/services/auth"
	"github.com/Faner201/go-grpc-server/sso/internal/storage/postrgres"
	"github.com/rs/zerolog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *zerolog.Logger, grpcPort int, toketTTL time.Duration, cfg *config.DatabaseConfig) *App {
	storage, err := postrgres.New(
		cfg.Host, cfg.User, cfg.Password, cfg.DBname, cfg.SSLmode, cfg.TimeZone, cfg.Port,
	)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, toketTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
