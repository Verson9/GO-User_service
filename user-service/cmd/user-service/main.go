package main

import (
	"GO-User_service/user-service/app"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app.Run()
}
