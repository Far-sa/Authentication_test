package migrator

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	MigrateUp() error
	MigrateDown() error
}

// TODO: clarify that exported or unexported properties
type MysqlMigrator struct {
	Db *sqlx.DB
	//config    ports.Config
	MigrationsDir string
}

func New(db *sqlx.DB, migrationsDir string) *MysqlMigrator {
	return &MysqlMigrator{
		Db:            db,
		MigrationsDir: migrationsDir,
	}
}

func (m *MysqlMigrator) MigrateUp() error {

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

func (m *MysqlMigrator) MigrateDown() error {

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
