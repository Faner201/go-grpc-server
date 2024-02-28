package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Faner201/go-grpc-server/sso/internal/app"
	"github.com/Faner201/go-grpc-server/sso/internal/config"
	"github.com/Faner201/go-grpc-server/sso/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	application := app.New(&log, cfg.GRPC.Port, cfg.TokenTTL, &cfg.Database)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info().Msg("Gracefuly stopped")
}
