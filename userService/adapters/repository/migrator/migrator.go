package migrator

import (
	"database/sql"
	"fmt"
	"user-svc/ports"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s",
		m.config.GetDatabaseConfig().User, m.config.GetDatabaseConfig().Password,
		m.config.GetDatabaseConfig().Host, m.config.GetDatabaseConfig().Port,
		m.config.GetDatabaseConfig().DBName))

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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s",
		m.config.GetDatabaseConfig().User, m.config.GetDatabaseConfig().Password,
		m.config.GetDatabaseConfig().Host, m.config.GetDatabaseConfig().Port,
		m.config.GetDatabaseConfig().DBName))

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
