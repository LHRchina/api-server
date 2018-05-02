package model

import (
	"api-server/src/util"
	"errors"
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

	var key string
	if uid > oid {
		key = getConcurrentKey(uid, oid)
	} else {
		key = getConcurrentKey(oid, uid)
	}

	dbObj.getConnect()
	conn := getRedisConn()
	defer (*conn).Close()

	ret, err := (*conn).Do("INCR", key)
	if err != nil {
		util.Err("model.UpdateRelationships  reids err:", err)
		return relationship, err
	}
	_, err = (*conn).Do("EXPIRE", key, DURATION)
	if err != nil {
		util.Err("model.UpdateRelationships  set reids  key expire err :", err)
		return relationship, errors.New("set expire fail")
	}
	val := ret.(int64)
	if val != 1 {
		util.Err("model.UpdateRelationships  reids  lock :")
		return relationship, errors.New("lock concurrent")
	}

	//查询是否存在关系
	util.Info("model.UpdateRelationships uid oid relation update :", uid, oid)
	err = dbObj.db.Model(&relationship).Where("uid = ? and oid = ?", uid, oid).Select()
	if err != nil && err != pg.ErrNoRows {
		util.Err("select error:", err)
		return relationship, err
	}
	//relation 不存在 插入或者更新
	if err != nil && err == pg.ErrNoRows {
		relationship.Status = dbState
		err = dbObj.db.Update(&relationship)
		if err != nil {
			util.Err("model.UpdateRelationships Update err:", err)
			return relationship, err
		}
	} else {
		relationship.Uid = uid
		relationship.Oid = oid
		relationship.Type = "relationship"
		relationship.Status = dbState
		err = dbObj.db.Insert(&relationship)
		if err != nil {
			util.Err("model.UpdateRelationships insert err:", err)
			return relationship, err
		}
	}
	var orelation Relationship
	err = dbObj.db.Model(&orelation).Where("uid = ? and oid = ?", oid, uid).Select()
	if err != nil {
		util.Err("model.UpdateRelationships select oid err:", err)
		return relationship, err
	}

	if orelation.Status == relationship.Status && orelation.Status == RelationStrToCode[LIKED] {
		//更新数据库为match
		orelation.Status = MATCHED
		relationship.Status = MATCHED
		ret, err := dbObj.db.Model(&relationship).Where("id = ?", relationship.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail relationship:", relationship, err)
			return relationship, err
		}
		util.Info("model.UpdateRelationships update match relationship:", relationship, "ret:", ret)

		ret, err = dbObj.db.Model(&orelation).Where("id = ?", orelation.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail orelation:", orelation, err)
			return relationship, err
		}
		util.Info("model.UpdateRelationships update match relationship:", orelation, "ret:", ret)
	}
	//之前是匹配的，后来disliked了，match状态改为liked
	if relationship.Status == RelationStrToCode[DISLIKED] && orelation.Status == RelationStrToCode[MATCHED] {
		orelation.Status = LIKED
		_, err = dbObj.db.Model(&orelation).Where("id = ?", orelation.Id).Update()
		if err != nil {
			util.Err("model.UpdateRelationships update match fail orelation:", orelation, err)
			return relationship, err
		}
	}
	_, err = (*conn).Do("DEL", key)
	if err != nil {
		util.Warn("model.UpdateRelationships del key:", key, "err:", err)
	}

	return relationship, nil
}
