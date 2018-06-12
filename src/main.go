package main

import (
	"api-server/src/handler"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", handler.GetAllUsersHandler).Methods("GET")
	r.HandleFunc("/users", handler.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{user_id:[0-9]+}/relationships", handler.GetUserRelationHandler).Methods("GET")
	r.HandleFunc("/users/{user_id:[0-9]+}/relationships/{other_user_id:[0-9]+}", handler.UpdateUserRelationHandler).Methods("PUT")
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("listen and serve err:", err)
		os.Exit(-1)
	}
}
