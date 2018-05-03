package model

import (
	"api-server/src/util"
	"fmt"
	"github.com/go-pg/pg"
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

const LIKE_STR = "liked"
const DISLIKE_STR = "disliked"
const MATCHED_STR = "matched"

const CONCURRENT_KEY_PRE = "concurent_%d_%d"
const DURATION = 10

func getConcurrentKey(uid, oid int64) string {
	return fmt.Sprintf(CONCURRENT_KEY_PRE, uid, oid)
}

func GetRelationshipsById(id int64) (relationships []Relationship, err error) {
	dbObj.getConnect()
	//TODO the sample is limit to 100
	err = dbObj.db.Model(&relationships).Where("uid = ?", id).Limit(100).Select()
	if err != nil {
		util.Err("GetRelationshipsById db.Model:", err)
		return relationships, err
	}
	for key, v := range relationships {
		relationships[key].Status = RelationCodeToStr[v.Status]
	}
	return relationships, nil
}

func UpdateRelationships(uid, oid int64, state string) (relationship Relationship, err error) {
	dbState := RelationStrToCode[state]
	dbObj.getConnect()
	//开启事务
	var tx *pg.Tx = new(pg.Tx)
	tx, err = dbObj.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err != nil {
		util.Err("model.UpdateRelationships begin fail:", err)
		return relationship, err
	}
	//查询是否存在关系
	util.Info("model.UpdateRelationships uid:", uid, " oid:", oid, " relation update ")
	_, err = tx.Query(&relationship, "select id,uid,oid,status,type from Relationships where uid = ? and oid = ? for update", uid, oid)
	if err != nil && err != pg.ErrNoRows {
		util.Err("model.UpdateRelationships select err:", err)
		return relationship, err
	}
	//relation 不存在 插入或者更新
	if relationship.Id == 0 {
		relationship.Uid = uid
		relationship.Oid = oid
		relationship.Type = "relationship"
		relationship.Status = dbState
		err := tx.Insert(&relationship)
		if err != nil {
			util.Err("model.UpdateRelationships Insert err:", err)
			return relationship, err
		}
	} else {
		relationship.Uid = uid
		relationship.Oid = oid
		relationship.Type = "relationship"
		relationship.Status = dbState
		err = tx.Update(&relationship)
		if err != nil {
			util.Err("model.UpdateRelationships Update err:", err)
			return relationship, err
		}
	}
	var orelation Relationship
	_, err = tx.Query(&orelation, "select  id,uid,oid,status,type from Relationships where uid = ? and oid = ? for update", oid, uid)
	if err != nil && err != pg.ErrNoRows {
		util.Err("model.UpdateRelationships select oid err:", err)
		return relationship, err
	}

	if orelation.Status == relationship.Status && orelation.Status == RelationStrToCode[LIKE_STR] {
		//更新数据库为match
		orelation.Status = MATCHED
		relationship.Status = MATCHED
		ret, err := tx.Model(&relationship).Where("id = ?", relationship.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail relationship:", relationship, err)
			return relationship, err
		}
		util.Info("model.UpdateRelationships update match relationship:", relationship, "ret:", ret)

		ret, err = tx.Model(&orelation).Where("id = ?", orelation.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail orelation:", orelation, err)
			return relationship, err
		}
		util.Info("model.UpdateRelationships update match relationship:", orelation, "ret:", ret)
	}
	//之前是匹配的，后来disliked了，match状态改为liked
	if relationship.Status == RelationStrToCode[DISLIKE_STR] && orelation.Status == RelationStrToCode[MATCHED_STR] {
		orelation.Status = LIKED
		_, err = tx.Model(&orelation).Where("id = ?", orelation.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail orelation:", orelation, err)
			return relationship, err
		}
	}

	return relationship, nil
}
