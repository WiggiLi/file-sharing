package store

import (
	"fmt"

	"github.com/WiggiLi/file-sharing-api/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"

	//_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" 	//?
)

func runPgMigrations() error {
	cfg:= config.Get()
	if cfg.PgMigrationsPath == ""{
		return errors.New("No cfg.PgMigrationsPath provided")
	}
	if cfg.PgURL == "" {
		return errors.New("No cfg.PgURL provided")
	}
	m, err := migrate.New(
		cfg.PgMigrationsPath,
		cfg.PgURL,
	)
	if err != nil {
		fmt.Println("HERE11")
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("HERE22")
		return err
	}
	return nil
}