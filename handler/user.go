package handler

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"

	"api-server/model"
)

const GET_METHOD_STR = "GET"

var commonRespMsg map[string]string

type PostUserData struct {
	Name string `json:"name"`
}

type PutState struct {
	State string `json:"state"`
}

func init() {
	//init the return msg
	commonRespMsg = map[string]string{
		"method_error": `{"code":200,"msg":"request method error %s"}`,
		"params_error": `{"code":200,"msg":"request method error %s"}`,
	}
}

//
func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-type", "application/json")

	ret , err := model.GetAllUser()
	if err != nil {
		//TODO
	}
	resultData , err := json.Marshal(ret)
	if err != nil {
		//TODO
	}
	fmt.Fprintf(w,string(resultData))
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	defer r.Body.Close()

	rawData,err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("http request body can't read:",err)
		fmt.Fprintf(w,"wrong")
	}

	var user PostUserData

	err = json.Unmarshal(rawData,&user)
	if err != nil {
		log.Println("http post Data error:",err)
		fmt.Fprintf(w,"wrong")
	}

	if  user.Name == "" {
		log.Println("http post Data error:",err)
		fmt.Fprintf(w,"wrong")
	}
	ret,err := model.InsertUser(user.Name)
	if err != nil {
		log.Println("insert err:",err)
	}
	retRaw ,err := json.Marshal(ret)
	if err != nil {
		log.Println("insert marshal err:",err)
	}
	fmt.Fprintf(w, string(retRaw))
}
