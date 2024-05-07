package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	// Optional: Use sqlx for convenience methods
	_ "github.com/lib/pq" // Import postgres driver (assuming you're using Postgres)
)

type Config struct {
	User     string
	Password string
	Port     int
	Host     string
	DbName   string
}

var (
	once sync.Once
	pool *sql.DB // Use sqlx.DB for convenience methods (optional)
)

// GetConnectionPool establishes a connection pool or returns the existing one
func GetConnectionPool(cfg Config) (*sql.DB, error) {
	once.Do(func() {
		var err error
		dataSourceName := fmt.Sprintf("postgres://%s:%s@tcp(%s:%d)%s?sslmode=disable",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
		pool, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}
		// Optional connection pool configuration (e.g., pool size)
		pool.SetMaxOpenConns(10) // Set the maximum number of open connections
		pool.SetMaxIdleConns(5)  // Set the maximum number of idle connections
	})
	return pool, nil
}
