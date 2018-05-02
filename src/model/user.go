package model

import (
	"api-server/src/util"
)

const userType = "user"

func GetAllUser() (user []User, err error) {
	dbObj.getConnect()
	//TODO the sample is limit to 100
	err = dbObj.db.Model(&user).Limit(100).Select()
	if err != nil {
		util.Err("GetAllUser err:", err)
		return nil, err
	}
	return
}

func InsertUser(name string) (u User, err error) {
	dbObj.getConnect()
	//TODO 如果存在同名的用户是否允许？
	u.Name = name
	u.Type = userType

	err = dbObj.db.Insert(&u)
	if err != nil {
		util.Err("insert error")
		return u, err
	}
	return u, nil
}
