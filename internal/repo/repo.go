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

func Init(pool *pg.CustomSqlConn) func() {
	lock.Lock()
	defer lock.Unlock()
	if singleton == nil {
		singleton = &repo{Db: pool}
	}
	return destroy
}

func destroy() {
	singleton = nil
}

func GetInstance() *repo {
	return singleton
}

func ScanRow(rows *pg.CustomRows, pointers ...interface{}) error {
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			return err
		}
	}
	return nil
}
