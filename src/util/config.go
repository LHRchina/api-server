package util

import (
	"encoding/json"
	"log"
	"os"
)

type (
	Db struct {
		Host   string `json:"host"`
		Port   string `json:"port"`
		User   string `json:"user"`
		DbName string `json:"db_name"`
		Password string `json:"password"`
	}

	Config struct {
		Pgdb Db `json:"postgresql_server"`
		RedisDb Db `json:"redis"`
	}
)

var conf Config

func init() {
	fl, err := os.Open("etc/config.json")
	if err != nil {
		log.Println("open file:etc/config.json err:", err)
		os.Exit(-1)
	}
	defer fl.Close()
	decoder := json.NewDecoder(fl)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println("decode file:etc/config.json err:", err)
		os.Exit(-1)
	}
}

func GetDb() *Db {
	return &conf.Pgdb
}

func GetRedis() *Db {
	
}