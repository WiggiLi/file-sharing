package store

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/WiggiLi/file-sharing-api/store/local"
	"github.com/WiggiLi/file-sharing-api/lib/logger"
	"github.com/WiggiLi/file-sharing-api/config"
	"github.com/WiggiLi/file-sharing-api/store/pg"
	"github.com/WiggiLi/file-sharing-api/store/redis"
	"github.com/WiggiLi/file-sharing-api/model"
)
type UserRepo struct {
	Pg model.UserRepoPg
	Redis model.UserRepoRedis
}

//Store contains all repositories
type Store struct {
	Pg *pg.DB
	Redis *redis.DB

	User UserRepo
	File model.FileMetaRepo
	FileContent model.FileContentRepo 
}

func New(ctx context.Context) (*Store, error) {
	cfg := config.Get()

	//Connect to postgres
	resisDB, err := redis.Dial(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "redisdb.Dial failed")
	}

	//Connect to postgres
	pgDB, err := pg.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "pgdb.Dial failed")
	}

	// Run Postgres migrations
	if pgDB != nil {
		log.Println("Running PostgreSQL migrations...")
		if err := runPgMigrations(); err != nil {
			return nil, errors.Wrap(err, "runPgMigrations failed")
		}
	}

	var store Store

	// Init Postgres repositories
	if pgDB != nil {
		store.Pg = pgDB
		go store.KeepAlivePg()
		
		store.User.Pg = pg.NewUserRepo(pgDB)
		store.File = pg.NewFileMetaRepo(pgDB)
	}

	if resisDB != nil {
		store.Redis = resisDB
		store.User.Redis = redis.NewUserRepo(resisDB)	
	}
	if cfg.FilePath != ""{
		store.FileContent = local.NewFileContentRepo(cfg.FilePath)
	}
	return &store, nil
}

// KeepAlivePollPeriod is a Pgkeepalive check time period
const KeepAlivePollPeriod = 3

// KeepAlivePg makes sure PostgreSQL is alive and reconnects if needed
func (store *Store) KeepAlivePg() {
	logger := logger.Get()
	var err error
	for {
		// Check if PostgreSQL is alive every 3 seconds
		time.Sleep(time.Second * KeepAlivePollPeriod)
		lostConnect := false
		if store.Pg == nil {
			lostConnect = true
		} else if _, err = store.Pg.Exec("SELECT 1"); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		logger.Debug().Msg("[store.KeepAlivePg] Lost PostgreSQL connection. Restoring...")
		store.Pg, err = pg.Dial()
		if err != nil {
			logger.Err(err)
			continue
		}
		logger.Debug().Msg("[store.KeepAlivePg] PostgreSQL reconnected")
	}
}
