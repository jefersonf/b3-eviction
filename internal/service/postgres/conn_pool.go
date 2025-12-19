package postgres

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewConnectionPool() (*pgxpool.Pool, error) {

	cfg := Config{
		Host:     cmp.Or(os.Getenv("POSTGRES_HOST"), "postgres"),
		Port:     cmp.Or(os.Getenv("POSTGRES_PORT"), "5432"),
		User:     cmp.Or(os.Getenv("POSTGRES_USER"), "admin"),
		Password: cmp.Or(os.Getenv("POSTGRES_PASSWORD"), "secret"),
		DBName:   cmp.Or(os.Getenv("POSTGRES_DB"), "votings"),
	}

	dbPool, err := NewConnection(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return dbPool, nil
}

func NewConnection(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	var err error

	once.Do(func() {

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)

		poolConfig, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			err = fmt.Errorf("unable to parse database config: %w", parseErr)
			return
		}

		// Optimization: Set pool settings for high throughput
		poolConfig.MaxConns = 20                   // Adjust based on your CPU/Postgres limits
		poolConfig.MinConns = 5                    // Keep some connections warm
		poolConfig.MaxConnLifetime = 1 * time.Hour // Recycle connections periodically
		poolConfig.MaxConnIdleTime = 30 * time.Minute

		// Create the pool
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		// Verify connection immediately
		if pingErr := pool.Ping(ctx); pingErr != nil {
			err = fmt.Errorf("database unreachable: %w", pingErr)
			return
		}

		log.Println("Successfully connected to Postgres (pgxpool)")
	})

	return pool, err
}

// Shutdown closes the pool
func Shutdown(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
