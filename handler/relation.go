package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"api-server/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func GetUserRelationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	defer r.Body.Close()

	uidStr := mux.Vars(r)["user_id"]
	uid , err:= strconv.Atoi(uidStr)
	if err != nil {
		log.Println("[GetUserRelationHandler] uid covert to int err:",err)
	}
	//判断uid是否正常
	if uid == 0 {
		fmt.Fprintf(w,"wrong")
	}

	ret := model.GetRelationshipsById(int64(uid))
	retJson,_ := json.Marshal(ret)
	fmt.Fprintf(w, string(retJson))
}

func UpdateUserRelationHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	defer r.Body.Close()
	rawData,err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("err:",err)
	}
	var state PutState
	err = json.Unmarshal(rawData, &state)
	if err != nil {
		log.Println("request err")
		fmt.Fprintf(w, "wrong")
		return
	}

	if state.State == ""{
		fmt.Fprintf(w,"wrong")
		return
	}
	uidStr := mux.Vars(r)["user_id"]
	oidStr := mux.Vars(r)["other_user_id"]
	//TODO 判断uid是否正常
	if uidStr == "0" || oidStr == "0"{
		log.Println("uid is fail")
		fmt.Fprintf(w,"wrong")
		return
	}

	uid,_ := strconv.Atoi(uidStr)
	oid,_ := strconv.Atoi(oidStr)
	ret , err := model.UpdateRelationships(int64(uid),int64(oid),state.State)
	retRaw,_ := json.Marshal(ret)
	fmt.Fprintf(w, string(retRaw))
}
