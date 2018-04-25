package model

import (
	"log"
	"fmt"
)

var relationMap = map[string]string {
	"0":"liked",
	"1":"disliked",
	"2":"matched",
}

var RelationcodeToStr = map[string]string{
	"liked":"0",
	"disliked":"1",
	"matched":"2",
}


func GetRelationshipsById(id int64) (relationships []Relationship ){
	db := getDbConnect()
	defer db.Close()

	//TODO add limit
	err := db.Model(&relationships).Where("uid = ?", id).Select()
	if err != nil {
		log.Println("db err:",err)
		return nil
	}
	for key,v := range relationships {
		relationships[key].Status = relationMap[v.Status]
	}
	return relationships
}

func UpdateRelationships(uid, oid int64, state string) (relationship Relationship  ,err error){
	//查询是否存在关系
	db := getDbConnect()
	defer db.Close()
	dbState := RelationcodeToStr[state]
	//TODO add limit
	log.Println(uid,oid)
	err = db.Model(&relationship).Where("uid = ? and oid = ?", uid, oid).Select()
	if err != nil {
		log.Println("select error:",err)
		//return relationship, err
	}
	//log.Println("ret:",relationship)
	//TODO relation存在
	if relationship.Type != "" {
		fmt.Println("exist")
		return relationship,nil
	}
	//不存在
	relationship.Uid = uid
	relationship.Oid = oid
	relationship.Type = "relationship"
	relationship.Status = dbState
	err = db.Insert(&relationship)
	if err != nil {
		log.Println("insert err:",err)
		return relationship, nil
	}

	var orelation Relationship
	err = db.Model(&orelation).Where("uid = ? and oid = ?",oid,uid).Select()
	if err != nil {
		return relationship,err
	}

	if orelation.Status == relationship.Status && orelation.Status == "0" {
		//更新数据库为match
		orelation.Status = MATCHED
		relationship.Status = MATCHED
		ret,err := db.Model(&relationship).Where("id = ?", relationship.Id).Update()
		if err != nil {
			log.Println("1",err)
		}
		log.Println(ret)
		ret,err = db.Model(&orelation).Where("id = ?", orelation.Id).Update()
		if err != nil {
			log.Println("2",err)
		}
		log.Println(ret)
		return relationship,nil
	}


	return relationship,nil
}