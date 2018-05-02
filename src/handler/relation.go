package handler

import (
	"api-server/src/model"
	"api-server/src/util"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

//get all relation of request user
func GetUserRelationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	defer r.Body.Close()

	//log the raw header and body info
	rawData, _ := ioutil.ReadAll(r.Body)
	util.Info("body :", rawData, "header:", r.Header)

	uidStr := mux.Vars(r)["user_id"]
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		util.Warn("handler.GetUserRelationHandler uid covert to int err:", err)
		fmt.Fprintf(w, string(dummySlice))
		return
	}

	//判断uid是否正常
	if uid == 0 {
		util.Warn("handler.GetUserRelationHandler uid is 0")
		fmt.Fprintf(w, string(dummySlice))
		return
	}

	ret, err := model.GetRelationshipsById(int64(uid))
	if err != nil {
		util.Warn("GetUserRelationHandler model.GetRelationshipsById uid:", uid, "err:", err)
		fmt.Fprintf(w, string(dummySlice))
		return
	}
	retJson, err := json.Marshal(ret)
	if err != nil {
		fmt.Fprintf(w, string(dummySlice))
		return
	}

	fmt.Fprintf(w, string(retJson))
}

//update relation ship
func UpdateUserRelationHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	defer r.Body.Close()
	rawData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Info("handler.UpdateUserRelationHandler ioutil.ReadAll err:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	//log the raw header and body info
	util.Info("body :", string(rawData), "header:", r.Header)

	var state PutState
	err = json.Unmarshal(rawData, &state)
	if err != nil {
		util.Info("handler.UpdateUserRelationHandler json.Unmarshal err:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	_, ok := model.StateMap[state.State]
	if !ok {
		util.Warn("handler.UpdateUserRelationHandler params invalid")
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	uidStr := mux.Vars(r)["user_id"]
	oidStr := mux.Vars(r)["other_user_id"]
	//check uid is valid
	if uidStr == "0" || oidStr == "0" {
		util.Warn("handler.UpdateUserRelationHandler params invalid")
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	uid, _ := strconv.Atoi(uidStr)
	oid, _ := strconv.Atoi(oidStr)
	ret, err := model.UpdateRelationships(int64(uid), int64(oid), state.State)
	if err != nil {
		util.Warn("handler.UpdateUserRelationHandler model fail:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}
	retRaw, _ := json.Marshal(ret)
	fmt.Fprintf(w, string(retRaw))
}
