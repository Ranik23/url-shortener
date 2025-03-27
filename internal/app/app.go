package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ranik23/url-shortener/api/proto/gen"
	"github.com/Ranik23/url-shortener/internal/config"
	grpc_server "github.com/Ranik23/url-shortener/internal/controllers/grpc"
	http_controllers "github.com/Ranik23/url-shortener/internal/controllers/http"
	"github.com/Ranik23/url-shortener/internal/libs/closer"
	http_server "github.com/Ranik23/url-shortener/internal/libs/http_server"
	repo_helpers "github.com/Ranik23/url-shortener/internal/libs/repository_helpers"
	"github.com/Ranik23/url-shortener/internal/repository/postgres"
	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
)


type App struct {
	config					*config.Config
	HTTPshortenerServer 	*http_server.Server
	gRPCshortenerServer 	*grpc.Server
	closer					*closer.Closer
	logger					*slog.Logger
}

func NewApp() (*App, error) {

	logger 	 := slog.New(tint.NewHandler(os.Stdout, nil))
	closer 	 := closer.NewCloser()
	cfg, err := config.LoadConfig(".env", "config/config.yaml")
	if err != nil {
		return nil, err
	}

	connectionString := repo_helpers.GetConnectionString(cfg.Database.Type, 
		cfg.Database.Host, cfg.Database.Port, cfg.Database.UserName, cfg.Database.Password, cfg.Database.DBName)
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	closer.Add(func(_ context.Context) error {
		pool.Close()
		return nil
	})
	
	txManager 	:= postgres.NewTxManager(pool, slog.Default())
	linkRepo 	:= postgres.NewPostgresLinkRepository(txManager)
	userRepo 	:= postgres.NewPostgresUserRepository(txManager)
	linkService := service.NewLinkService(linkRepo, txManager)
	statService := service.NewStatService()
	userService := service.NewUserService(userRepo, txManager)
	service 	:= service.NewService(linkService, statService, userService)
	handler 	:= http_controllers.NewHandler(service)

	handler.SetUpRoutes()
	
	grpcServer 			:= grpc.NewServer(grpc.UnaryInterceptor(grpc_server.ErrorHandlerInterceptor))
	grpcShortenerServer := grpc_server.NewShortenerServer(service)

	gen.RegisterURLShortenerServer(grpcServer, grpcShortenerServer)

	closer.Add(func(_ context.Context) error {
		grpcServer.GracefulStop()
		return nil
	})

	httpConfig := http_server.Config{
		Host: cfg.HTTPServer.Host,
		Port: cfg.HTTPServer.Port,
		StartMsg: "Hello",
		ShutdownTimeout: 10 * time.Second,
	}
	shortenerServer := http_server.New(logger, httpConfig, handler)

	return &App{
		HTTPshortenerServer: shortenerServer,
		gRPCshortenerServer: grpcServer,
		config: cfg,
		closer: closer,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		if err := a.closer.Close(ctx); err != nil {
			a.logger.Error("Failed to close resources", slog.String("error", err.Error()))
		}
		a.logger.Info("Successfully closed all resources")
	}()

	grpcAddr 	  := fmt.Sprintf("%s:%s", a.config.GRPCServer.Host, a.config.GRPCServer.Port)
	errorCh 	  := make(chan error, 2)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}	

	go func() {
		if err := a.HTTPshortenerServer.Start(context.Background()); err != nil {
			errorCh <- fmt.Errorf("failed to start gRPC Server: %v", err)
		}
	}()

	go func() {
		if err := a.gRPCshortenerServer.Serve(listener); err != nil {
			errorCh <- fmt.Errorf("failed to start HTTP Server: %v", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)


	select {
	case <- quit:
		a.logger.Info("Получен сигнал завершения, выключаем gRPC сервер...")
		return nil
	case err := <- errorCh:
		a.logger.Error("Ошибка сервера", slog.String("error", err.Error()))
		return err
	}
}