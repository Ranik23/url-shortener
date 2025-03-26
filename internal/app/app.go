package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/Ranik23/url-shortener/api/proto/gen"
	"github.com/Ranik23/url-shortener/internal/config"
	grpcserver "github.com/Ranik23/url-shortener/internal/controllers/grpc"
	"github.com/Ranik23/url-shortener/internal/controllers/http"
	"github.com/Ranik23/url-shortener/internal/repository/postgres"
	"github.com/Ranik23/url-shortener/internal/server"
	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
)


type App struct {
	config					*config.Config
	HTTPshortenerServer 	*server.ShortenerServer
	gRPCshortenerServer 	*grpc.Server
}

func NewApp() (*App, error) {

	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	
	closer := NewCloser()

	cfg, err := config.LoadConfig(".env", "config/config.yaml")
	if err != nil {
		return nil, err
	}

	dsn := cfg.Database.GetPostgresDSN()
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	closer.Add(func(ctx context.Context) error {
		pool.Close()
		return nil
	})
	

	txManager := postgres.NewTxManager(pool, slog.Default())
	linkRepo := postgres.NewPostgresLinkRepository(txManager)
	userRepo := postgres.NewPostgresUserRepository(txManager)
	linkService := service.NewLinkService(linkRepo, txManager)
	statService := service.NewStatService()
	userService := service.NewUserService(userRepo, txManager)
	service := service.NewService(linkService, statService, userService)
	handler := http.NewHandler(service)

	handler.SetUpRoutes()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcserver.ErrorHandlerInterceptor))
	grpcShortenerServer := grpcserver.NewShortenerServer(service)

	gen.RegisterURLShortenerServer(grpcServer, grpcShortenerServer)


	shortenerServer := server.NewShortenerServer("localhost:8080", handler, logger)

	closer.Add(func(ctx context.Context) error {
		if err := shortenerServer.ShutDown(ctx); err != nil {
			return err
		}
		return nil
	})

	return &App{
		HTTPshortenerServer: shortenerServer,
		gRPCshortenerServer: grpcServer,
		config: cfg,
	}, nil
}

func (a *App) Run() {

	go func() {
		a.HTTPshortenerServer.ListenServe()
	}()

	go func() {

		address := fmt.Sprintf("%s:%s", a.config.Database.Host, a.config.Database.Port)
		listener, err := net.Listen("tcp", address)

		a.gRPCshortenerServer.Serve()
	}()
}