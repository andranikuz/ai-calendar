package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

type Config struct {
	Host         string
	Port         int
	Database     string
	User         string
	Password     string
	SSLMode      string
	MaxConns     int
	MinConns     int
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
}

func NewDatabase(config Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Database,
		config.User,
		config.Password,
		config.SSLMode,
	)
	
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}
	
	// Configure connection pool
	pgxConfig.MaxConns = int32(config.MaxConns)
	pgxConfig.MinConns = int32(config.MinConns)
	pgxConfig.MaxConnLifetime = config.MaxLifetime
	pgxConfig.MaxConnIdleTime = config.MaxIdleTime
	
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &Database{
		Pool: pool,
	}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *Database) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// GetDefaultConfig returns default database configuration
func GetDefaultConfig() Config {
	return Config{
		Host:         "localhost",
		Port:         5432,
		Database:     "smart_calendar",
		User:         "postgres",
		Password:     "postgres",
		SSLMode:      "disable",
		MaxConns:     10,
		MinConns:     2,
		MaxLifetime:  time.Hour,
		MaxIdleTime:  time.Minute * 30,
	}
}