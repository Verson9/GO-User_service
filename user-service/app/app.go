package app

import (
	"GO-User_service/user-service/api"
	"GO-User_service/user-service/internal/handler"
	"GO-User_service/user-service/internal/usersdb"
	"fmt"
	"log"
	"net/http"
)

const Port = "8080"

func Run() {
	u := usersdb.Connect()
	l := log.Logger{}

	s := handler.NewHandler(&u, &l)
	h := api.Handler(s)

	err := http.ListenAndServe(Port, h)
	if err != nil {
		l.Fatalln(fmt.Sprintf("failed to serve connection on port \"%s\"", Port))
	}
}
