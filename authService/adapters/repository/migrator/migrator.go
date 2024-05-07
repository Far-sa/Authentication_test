package migrator

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	MigrateUp() error
	MigrateDown() error
}

type PgMigrator struct {
	//config        PgConfig
	DB            *sql.DB // Use connection pool
	MigrationsDir string
	//	driverURL     string // Store the constructed driver URL here

}

func NewMigrator(db *sql.DB, migrationsDir string) *PgMigrator {
	return &PgMigrator{
		DB:            db,
		MigrationsDir: migrationsDir,
	}

}

// ! The current implementation creates a new driverURL for each migration execution
// !(both MigrateUp and MigrateDown). While it still functions, it's slightly less efficient.
// driverURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
// 	m.config.User, m.config.Password, m.config.Host, m.config.Port, m.config.DbName)

func (m *PgMigrator) MigrateUp() error {

	migrator, err := migrate.New(
		"file://"+m.MigrationsDir, "")
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Println("No migrations found to apply (Up)")
	} else {
		log.Printf("Up migrations applied successfully!")
	}
	return nil
}

func (m *PgMigrator) MigrateDown() error {

	migrator, err := migrate.New(
		"file://"+m.MigrationsDir, "")
	if err != nil {
		return err
	}

	if err := migrator.Down(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Println("No migrations found to apply (Down)")
	} else {
		log.Printf("Down migrations applied successfully!")
	}
	return nil
}
