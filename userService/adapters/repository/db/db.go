package db

import (
	"fmt"
	"log"
	"sync"
	"user-svc/ports"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// type Config struct {
// 	config ports.Config
// }

var (
	once sync.Once
	pool *sqlx.DB // Use sqlx.DB for convenience methods (optional)
)

// GetConnectionPool establishes a connection pool or returns the existing one
func GetConnectionPool(cfg ports.Config) (*sqlx.DB, error) {

	dbConfig := cfg.GetDatabaseConfig()

	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

		fmt.Println("dsn :", dsn)

		pool, err = sqlx.Open("mysql", dsn)
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}

		if err := pool.Ping(); err != nil {
			log.Printf("Error pinging database: %v", err)
			// Additional error handling or cleanup can be added here if needed
		}

		if pool == nil {
			log.Fatal("Database object is nil")
		}

		// Optional connection pool configuration (e.g., pool size)
		pool.SetMaxOpenConns(10) // Set the maximum number of open connections
		pool.SetMaxIdleConns(5)  // Set the maximum number of idle connections
	})
	return pool, nil
}
