package repo

import (
	"jwt/pkg/helpers/pg"
	"sync"
)

var lock = &sync.Mutex{}

type repo struct {
	Db *pg.CustomSqlConn
}

var singleton *repo

func Init(pool *pg.CustomSqlConn) {
	lock.Lock()
	defer lock.Unlock()
	if singleton == nil {
		singleton = &repo{Db: pool}
	}
}

func getInstance() *repo {
	return singleton
}
