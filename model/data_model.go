package model

import (
	"github.com/go-pg/pg"
	"api-server/config"
)

type User struct {
	Id int64 `sql:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Relationship struct {
	Id int64 `json:"id"`
	Uid int64 `json:"uid"`
	Oid int64 `json:"oid"`
	Status string `json:"status"`
	Type string `json:"type"`
}

const (
	USER_TYPE = iota
	RELATIONSHIPS_TYPE
)

const (
	LIKED_ = iota
	DISLIKED_
	MATCHED_
)

const LIKED  = "0"
const DISLIKED = "1"
const MATCHED = "2"

var stateMap map[string]int = map[string]int{
	"liked":LIKED_,
	"disliked":DISLIKED_,
}

//TODO get conn from connect pool
func getDbConnect()  *pg.DB {
	db := config.GetDb()
	return pg.Connect(&pg.Options{
		Addr: db.Host+":"+db.Port,
		User: db.User,
		Database:db.DbName,
	})
}



