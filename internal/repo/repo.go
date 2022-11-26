package repo

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var lock = &sync.Mutex{}

type repo struct {
	Db *pgxpool.Pool
}

var singleton *repo

func Init(pool *pgxpool.Pool) {
	lock.Lock()
	defer lock.Unlock()
	if singleton == nil {
		singleton = &repo{Db: pool}
	}
}

func getInstance() *repo {
	return singleton
}
