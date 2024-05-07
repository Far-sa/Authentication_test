package db

import (
	"fmt"
	"log"
	"sync"
	"user-svc/ports"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	config ports.Config
}

var (
	once sync.Once
	pool *sqlx.DB // Use sqlx.DB for convenience methods (optional)
)

// GetConnectionPool establishes a connection pool or returns the existing one
func GetConnectionPool(cfg ports.Config) (*sqlx.DB, error) {

	dbConfig := cfg.GetDatabaseConfig()

	once.Do(func() {
		var err error
		dataSourceName := fmt.Sprintf("postgres://%s:%s@tcp(%s:%d)%s?sslmode=disable",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
		pool, err = sqlx.Open("postgres", dataSourceName)
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}
		// Optional connection pool configuration (e.g., pool size)
		pool.SetMaxOpenConns(10) // Set the maximum number of open connections
		pool.SetMaxIdleConns(5)  // Set the maximum number of idle connections
	})
	return pool, nil
}
