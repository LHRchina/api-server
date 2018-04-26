package model

import (
	"api-server/src/util"
	"github.com/go-pg/pg"
	"time"
)

//连接池初始化
var db *pg.DB

func init() {
	db = getDbConnect()
}

func getDbConnect() *pg.DB {
	db := util.GetDb()
	return pg.Connect(&pg.Options{
		Addr:               db.Host + ":" + db.Port,
		User:               db.User,
		Database:           db.DbName,
		PoolSize:           1000,
		IdleCheckFrequency: time.Second * 4,
	})
}
