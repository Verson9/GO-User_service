package app

import (
	"GO-User_service/user-service/api"
	"GO-User_service/user-service/internal/handler"
	"GO-User_service/user-service/internal/usersdb"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Run() {
	port := os.Getenv("APP_PORT")
	u := usersdb.Connect()
	l := log.Logger{}

	s := handler.NewHandler(&u, &l)
	h := api.Handler(s)

	err := http.ListenAndServe(":"+port, h)
	if err != nil {
		l.Fatalf(fmt.Sprintf("failed to serve connection on port \"%s\": %v", port, err))
	}
}
