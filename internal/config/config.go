package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

type HTTPServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

type GRPCServerConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

type PoolConfig struct {
	MaxConnections    int `yaml:"max_connections"`
	MinConnections    int `yaml:"min_connections"`
	MaxLifeTime       int `yaml:"max_lifetime"`
	MaxIdleTime       int `yaml:"max_idle_time"`
	HealthCheckPeriod int `yaml:"health_check_period"`
}

type DatabaseConfig struct {
	Type               string     `yaml:"type"`
	Host               string     `yaml:"host"`
	Port               int        `yaml:"port"`
	DBName             string     `yaml:"dbname"`
	SSL                string     `yaml:"ssl"`
	ConnectionAttempts int        `yaml:"connection_attempts"`
	Pool               PoolConfig `yaml:"pool"`
	// Значения для Username и Password загружаются из окружения
	Username string `yaml:"-"`
	Password string `yaml:"-"`
}

type CacheConfig struct {
	Type               string `yaml:"type"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Password           string `yaml:"password"`
	DB                 int    `yaml:"db"`
	ConnectionAttempts int    `yaml:"connection_attempts"`
}

type StorageConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Cache    CacheConfig    `yaml:"cache"`
}

type Config struct {
	HTTPServer HTTPServerConfig `yaml:"http_server"`
	GRPCServer GRPCServerConfig `yaml:"grpc_server"`
	Storage    StorageConfig    `yaml:"storage"`
}

func LoadConfig(envPath string, configPath string) (*Config, error) {

	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("unmarshalling YAML: %v", err)
	}

	config.Storage.Database.Host = os.Getenv("DB_HOST")
	config.Storage.Database.Username = os.Getenv("DB_USER_NAME")
	config.Storage.Database.Password = os.Getenv("DB_PASSWORD")

	config.Storage.Cache.Host = os.Getenv("REDIS_HOST")

	return &config, nil
}

func (s *StorageConfig) ConnectionToPostgres() (*pgxpool.Pool, error) {
	cfg := s.Database
	dsn := s.GetDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres DSN: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.Pool.MaxConnections)
	poolConfig.MinConns = int32(cfg.Pool.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.Pool.MaxLifeTime) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(cfg.Pool.MaxIdleTime) * time.Second
	poolConfig.HealthCheckPeriod = time.Duration(cfg.Pool.HealthCheckPeriod) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	var attempt int
	for attempt < cfg.ConnectionAttempts {
		if err := pool.Ping(context.Background()); err == nil {
			break
		}
		attempt++
		time.Sleep(2 * time.Second)
		log.Printf("Attempt to connect to PostgreSQL: attempt=%d, max_attempts=%d, error=%v",
			attempt, cfg.ConnectionAttempts, err)
	}

	if attempt == cfg.ConnectionAttempts {
		return nil, fmt.Errorf("connect to PostgreSQL after %d attempts", cfg.ConnectionAttempts)
	}

	return pool, nil
}

func (s *StorageConfig) GetDSN() string {
	db := s.Database
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.Username, db.Password, db.DBName, db.SSL,
	)
}

func (s *StorageConfig) ConnectionToRedis() (*redis.Client, error) {
	cache := s.Cache
	addr := fmt.Sprintf("%s:%d", cache.Host, cache.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cache.Password,
		DB:       cache.DB,
	})

	attempts := cache.ConnectionAttempts

	var attempt int
	var err error

	for attempt < attempts {
		err = client.Ping(context.Background()).Err()
		if err == nil {
			break
		}
		attempt++
		log.Printf("Attempt to connect to Redis: attempt=%d, max_attempts=%d, error=%v",
			attempt, attempts, err)
		time.Sleep(2 * time.Second)
	}

	if attempt == attempts {
		return nil, fmt.Errorf("connect to Redis after %d attempts: %w", attempts, err)
	}

	return client, nil
}
