package pg

import (
	"time"
	"errors"
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"

	"github.com/WiggiLi/file-sharing-api/config"
)

// Timeout is a Postgres timeout
const Timeout = 5

// DB is a shortcut structure to a Postgres DB
type DB struct {
	*pg.DB
}

type dbLogger struct { }

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
    return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) (error) {
	fq, _ := q.FormattedQuery()
    fmt.Println(string(fq))
    return nil
}

// Dial creates new database connection to postgres
func Dial() (*DB, error) {
	cfg:= config.Get()
	if cfg.PgURL == "" {
		return nil, errors.New("No cfg.PgURL provided")
	}
	pgOpts, err := pg.ParseURL(cfg.PgURL)
	if err != nil {
		return nil, err
	}

	pgDB := pg.Connect(pgOpts)

	pgDB.AddQueryHook(dbLogger{})

	_, err = pgDB.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	pgDB.WithTimeout(time.Second * time.Duration(Timeout))

	return &DB{pgDB}, nil
}