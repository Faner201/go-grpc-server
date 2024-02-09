package main

import (
	"github.com/Faner201/go-grpc-server/sso/internal/config"
	"github.com/Faner201/go-grpc-server/sso/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Println("Hello world")
}
