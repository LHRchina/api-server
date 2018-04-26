package model

import (
	"api-server/src/util"
)

var RelationCodeToStr = map[string]string{
	"0": "liked",
	"1": "disliked",
	"2": "matched",
}

var RelationStrToCode = map[string]string{
	"liked":    "0",
	"disliked": "1",
	"matched":  "2",
}

func GetRelationshipsById(id int64) (relationships []Relationship) {

	//TODO the sample is limit to 100
	err := db.Model(&relationships).Where("uid = ?", id).Limit(100).Select()
	if err != nil {
		util.Err("db err:", err)
		return
	}
	for key, v := range relationships {
		relationships[key].Status = RelationCodeToStr[v.Status]
	}
	return relationships
}

func UpdateRelationships(uid, oid int64, state string) (relationship Relationship, err error) {
	dbState := RelationStrToCode[state]

	//查询是否存在关系
	util.Info("model.UpdateRelationships uid oid relation update :",uid, oid)
	err = db.Model(&relationship).Where("uid = ? and oid = ?", uid, oid).Select()
	if err != nil {
		util.Err("select error:", err)
		return relationship,err
	}

	//TODO relation存在
	if relationship.Type != "" {
		util.Warn("model.UpdateRelationships relation exit uid:",uid,"oid:",oid)
		return relationship, nil
	}

	//不存在
	relationship.Uid = uid
	relationship.Oid = oid
	relationship.Type = "relationship"
	relationship.Status = dbState
	err = db.Insert(&relationship)
	if err != nil {
		util.Err("model.UpdateRelationships insert err:", err)
		return relationship, err
	}

	var orelation Relationship
	err = db.Model(&orelation).Where("uid = ? and oid = ?", oid, uid).Select()
	if err != nil {
		util.Err("model.UpdateRelationships select oid err:", err)
		return relationship, err
	}

	if orelation.Status == relationship.Status && orelation.Status == "0" {
		//更新数据库为match
		orelation.Status = MATCHED
		relationship.Status = MATCHED
		ret, err := db.Model(&relationship).Where("id = ?", relationship.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail relationship:",relationship, err)
		}
		util.Err("model.UpdateRelationships update match relationship:",relationship,"ret:",ret)

		ret, err = db.Model(&orelation).Where("id = ?", orelation.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail orelation:",orelation, err)
		}
		util.Err("model.UpdateRelationships update match relationship:",orelation,"ret:",ret)
		return relationship, nil
	}

	return relationship, nil
}
