package model

import (
	"api-server/src/util"
	"github.com/go-pg/pg"
	"sync"
	"time"
)

type DbObj struct {
	db *pg.DB
}

var dbObj DbObj
var once sync.Once

func init() {

}

func (d *DbObj) getConnect() {
	if d.db != nil {
		return
	}

	once.Do(func() {
		dbConf := util.GetDb()
		op := pg.Options{
			Addr:               dbConf.Host + ":" + dbConf.Port,
			User:               dbConf.User,
			Database:           dbConf.DbName,
			PoolSize:           1000,
			IdleCheckFrequency: time.Second * 4,
		}
		if dbConf.Password != "" {
			op.Password = dbConf.Password
		}
		d.db = pg.Connect(&op)
	})
}
