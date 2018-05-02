package model

import (
	"api-server/src/util"
	"github.com/go-pg/pg"
	redis2 "github.com/gomodule/redigo/redis"
	"os"
	"time"
)

type DbObj struct {
	db *pg.DB
}

var dbObj DbObj
var pool *redis2.Pool

func init() {
	redisDb := util.GetRedis()
	pool = &redis2.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis2.Conn, error) {
			add := redisDb.Host + ":" + redisDb.Port
			c, err := redis2.Dial("tcp", add)
			if err != nil {
				util.Err("init redis err:", err)
				os.Exit(-1)
			}
			if redisDb.Password != "" {
				if _, err := c.Do("AUTH", redisDb.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if redisDb.DbName != "" {
				if _, err := c.Do("SELECT", 0); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		MaxActive: 100,
	}
}

func getRedisConn() *redis2.Conn {
	conn := pool.Get()
	return &conn
}

func (d *DbObj) getConnect() {

	if d.db != nil {
		return
	}

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
}
