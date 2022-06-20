package app

import (
	han "GO-User_service/internal/handler"
	"GO-User_service/internal/usersdb"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() {
	users := usersdb.Connect()
	handler := han.NewHandler(&users)
	router := mux.NewRouter()

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HELLO!")
	})
	router.HandleFunc("/users", handler.GetUsers)
	router.HandleFunc("/user", handler.PostUser)
	router.HandleFunc("/user/{username}", handler.GetUsersByUsername)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
