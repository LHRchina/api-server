package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"api-server/src/model"
	"api-server/src/util"
)

var commonRespMsg map[string]string

var dummySlice, _ = json.Marshal([]interface{}{})
var dummyObj, _ = json.Marshal(struct{}{})

type PostUserData struct {
	Name string `json:"name"`
}

type PutState struct {
	State string `json:"state"`
}

//init the return msg
func init() {
	fmt.Println("init")

	commonRespMsg = map[string]string{
		"method_error": `{"code":200,"msg":"request method error %s"}`,
		"params_error": `{"code":200,"msg":"request method error %s"}`,
	}
}

//get all user list
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	//set response header Content-type json
	w.Header().Set("Content-type", "application/json")
	defer r.Body.Close()

	//log the raw header and body info
	rawData, _ := ioutil.ReadAll(r.Body)
	util.Info("body :", rawData, "header:", r.Header)

	//fetch all data
	ret, err := model.GetAllUser()
	if err != nil {
		util.Warn("handler.GetAllUsers err:", err)
		fmt.Fprintf(w, string(dummySlice))
		return
	}

	resultData, err := json.Marshal(ret)
	if err != nil {
		util.Warn("handler.GetAllUsers json.Marshal ret err", err)
		fmt.Fprintf(w, string(dummySlice))
		return
	}

	fmt.Fprintf(w, string(resultData))
}

//create user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	//set response header Content-type json
	w.Header().Set("Content-type", "application/json")
	defer r.Body.Close()

	rawData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Warn("http request body can't read:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	util.Info("body :", rawData, "header:", r.Header)

	var user PostUserData

	err = json.Unmarshal(rawData, &user)
	if err != nil {
		util.Warn("http post Data error:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	//user name should none empty
	if user.Name == "" {
		util.Warn("http post Data error:", err)
		fmt.Fprintf(w, "wrong")
		return
	}

	ret, err := model.InsertUser(user.Name)
	if err != nil {
		util.Warn("insert err:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}

	retRaw, err := json.Marshal(ret)
	if err != nil {
		util.Warn("insert marshal err:", err)
		fmt.Fprintf(w, string(dummyObj))
		return
	}
	fmt.Fprintf(w, string(retRaw))
}
