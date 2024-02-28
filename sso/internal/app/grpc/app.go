package appgrpc

import (
	"fmt"
	"net"
	"os"

	authgrpc "github.com/Faner201/go-grpc-server/sso/internal/grpc/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pereslava/grpc_zerolog"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *zerolog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zerolog.Logger, authService authgrpc.Auth, port int) *App {
	grpc_zerolog.ReplaceGrpcLogger(zerolog.New(os.Stderr).Level(zerolog.ErrorLevel))

	loggingOpts := []grpc_zerolog.Option{
		grpc_zerolog.WithLogOnEvents(
			grpc_zerolog.PayloadReceived, grpc_zerolog.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Panic().Any("Recovered from panic", p)
			return status.Errorf(codes.Internal, "Internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		grpc_zerolog.NewUnaryServerInterceptor(*log, loggingOpts...),
	))

	authgrpc.Register(gRPCServer, authService)
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
	const op = "gRPCapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info().Str("addr", l.Addr().String()).Msg("grpc server started")

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "gRPCapp.Stop"

	log := a.log.With().Str("op", op).Logger()
	log.Info().Int("port", a.port).Msg("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
