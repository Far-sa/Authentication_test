package migrator

import (
	"database/sql"
	"fmt"
	"user-svc/ports"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	//db        sqlx.DB
	config    ports.Config
	migration *migrate.FileMigrationSource
}

func New(config ports.Config) Migrator {

	migrations := &migrate.FileMigrationSource{
		Dir: "../../../infrastructure/db/migrations",
	}

	return Migrator{config: config, migration: migrations}
}

func (m Migrator) Up() error {

	dbConf := m.config.GetDatabaseConfig()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName))

	if err != nil {
		return err
	}

	_, eErr := migrate.Exec(db, "mysql", m.migration, migrate.Up)
	if eErr != nil {
		return eErr
	}

	return nil
}

func (m Migrator) Down() error {

	dbConf := m.config.GetDatabaseConfig()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName))

	if err != nil {
		return err
	}

	_, eErr := migrate.Exec(db, "mysql", m.migration, migrate.Down)
	if eErr != nil {
		return eErr
	}

	return nil
}

func (m Migrator) Status() {
	//TODO
}
