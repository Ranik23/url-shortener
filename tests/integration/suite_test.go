//go:build integration

package integration

import (
	"context"
	"database/sql"
	"github.com/Ranik23/url-shortener/api/proto/gen"
	"github.com/Ranik23/url-shortener/internal/config"
	grpc_server "github.com/Ranik23/url-shortener/internal/controllers/grpc"
	http_controllers "github.com/Ranik23/url-shortener/internal/controllers/http"
	"github.com/Ranik23/url-shortener/internal/repository"
	imppool "github.com/Ranik23/url-shortener/internal/repository/pgxpool"
	"github.com/Ranik23/url-shortener/internal/service"
	"github.com/Ranik23/url-shortener/tests/integration/testutil"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	psqlContainer *testutil.PostgreSQLContainer
	httpServer    *httptest.Server
	grpcServer    *grpc.Server

	svc       service.Service
	txManager repository.TxManager
	userRepo  repository.UserRepository
	linkRepo  repository.LinkRepository
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig("../../.env", "../../config/config.yaml")
	s.Require().NoError(err)

	psqlContainer, err := testutil.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	err = testutil.RunMigrations(psqlContainer.GetDSN(), "../../migrations")
	s.Require().NoError(err)

	poolConfig, err := pgxpool.ParseConfig(psqlContainer.GetDSN())
	s.Require().NoError(err)

	poolConfig.MaxConns = int32(cfg.Storage.Database.Pool.MaxConnections)
	poolConfig.MinConns = int32(cfg.Storage.Database.Pool.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.Storage.Database.Pool.MaxLifeTime)
	poolConfig.MaxConnIdleTime = time.Duration(cfg.Storage.Database.Pool.MaxIdleTime)
	poolConfig.HealthCheckPeriod = time.Duration(cfg.Storage.Database.Pool.HealthCheckPeriod)

	pgPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	s.Require().NoError(err)

	ctxManager := imppool.NewPgxCtxManager(pgPool)
	settings := imppool.NewPgxSettings()
	txManager := imppool.NewPgxTxManager(pgPool, logger, settings)

	linkRepo := imppool.NewPostgresLinkRepository(ctxManager, settings)
	userRepo := imppool.NewPgxUserRepository(ctxManager, settings)

	repo := repository.NewRepository(userRepo, linkRepo)

	linkService := service.NewLinkService(repo, txManager, logger)
	statService := service.NewStatService()
	userService := service.NewUserService(repo, txManager, logger)
	svc := service.NewService(linkService, statService, userService)

	s.svc = svc
	s.txManager = txManager
	s.userRepo = userRepo
	s.linkRepo = linkRepo

	handler := http_controllers.NewHandler(svc)
	router := gin.Default()
	handler.SetUpRoutes()
	s.httpServer = httptest.NewServer(router)

	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(grpc_server.ErrorHandlerInterceptor))
	grpcShortenerServer := grpc_server.NewShortenerServer(svc)
	gen.RegisterURLShortenerServer(s.grpcServer, grpcShortenerServer)
}

// Выполняется перед каждым тестом
func (s *TestSuite) SetupTest() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	_, err = db.Exec(`
        TRUNCATE TABLE users, links RESTART IDENTITY CASCADE;
    `)
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.httpServer.Close()
	s.grpcServer.GracefulStop()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
