package model

import (
	"log"
)

const userType  = "user"

func GetAllUser() (user []User, err error){
	db := getDbConnect()
	defer db.Close()

	//TODO add limit
	err = db.Model(&user).Select()
	if err != nil {
		log.Println("GetAllUser err:",err)
		return nil, err
	}
	return
}

func InsertUser(name string) (u User, err error) {
	db := getDbConnect()
	defer db.Close()

	//TODO 如果存在同名的用户是否允许？
	u.Name = name
	u.Type = userType
	err = db.Insert(&u)
	if err != nil {
		log.Println("insert error")
		return u,err
	}
	return u,nil
}