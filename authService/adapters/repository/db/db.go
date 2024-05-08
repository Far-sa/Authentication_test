package db

import (
	"auth-svc/internal/ports"
	"database/sql"
	"fmt"
	"log"
	"sync"

	// Optional: Use sqlx for convenience methods
	_ "github.com/lib/pq" // Import postgres driver (assuming you're using Postgres)
)

// type Config struct {
// 	config   ports.Config
// 	User     string
// 	Password string
// 	Port     int
// 	Host     string
// 	DbName   string
// }

var (
	once sync.Once
	pool *sql.DB // Use sqlx.DB for convenience methods (optional)
)

// GetConnectionPool establishes a connection pool or returns the existing one
func GetConnectionPool(cfg ports.Config) (*sql.DB, error) {

	dbCfg := cfg.GetDatabaseConfig()

	once.Do(func() {
		var err error
		dataSourceName := fmt.Sprintf("postgres://%s:%s@tcp(%s:%s)%s?sslmode=disable",
			dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)
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
