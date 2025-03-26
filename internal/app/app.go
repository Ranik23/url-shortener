package app

import (
	"context"
	"log/slog"

	"github.com/Ranik23/url-shortener/internal/config"
	"github.com/Ranik23/url-shortener/internal/controllers/http"
	"github.com/Ranik23/url-shortener/internal/repository/postgres"
	"github.com/Ranik23/url-shortener/internal/server"
	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)


type App struct {
	shortenerServer *server.ShortenerServer
}

func NewApp() (*App, error) {

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

	shortenerServer := server.NewShortenerServer("localhost:8080", handler, slog.Default())

	closer.Add(func(ctx context.Context) error {
		if err := shortenerServer.ShutDown(ctx); err != nil {
			return err
		}
		return nil
	})

	return &App{
		shortenerServer: shortenerServer,
	}, nil
}

// blocking
func (a *App) Run() {
	a.shortenerServer.ListenServe()
}