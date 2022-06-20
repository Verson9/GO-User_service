package usersdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type User struct {
	Username   string
	Password   string
	Privileges int
}

type Usersdb struct {
	db *sql.DB
}

func Connect() Usersdb {
	db := connectToDB()
	createUsersTableIfNotExist(db)
	return Usersdb{
		db: db,
	}
}

func connectToDB() *sql.DB {
	dbPassword := os.Getenv("DBPASS")
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(localhost:3306)/db", dbPassword))
	if err != nil {
		panic(err.Error())
	}
	log.Println("successfully connected to database.")
	return db
}

func createUsersTableIfNotExist(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, username VARCHAR(100) NOT NULL UNIQUE, password TEXT NOT NULL, privileges INT)")
	if err != nil {
		log.Fatalf("failed to create table \"users\": %v", err)
	}
	log.Println("successfully created table")
}

func (u *Usersdb) CreateUserIfNotExists(user User) error {
	_, err := u.db.Exec("INSERT INTO users ( username , password, privileges) VALUES (?, ?, ?)", user.Username, user.Password, user.Privileges)
	if err != nil {
		return fmt.Errorf("failed to create user %v, %v ", user, err)
	}
	log.Println("successfully created user")
	return nil
}

func (u *Usersdb) GetUser(username string) (User, error) {
	var user User
	data := u.db.QueryRow("SELECT username, password, privileges from users WHERE username = ?", username)
	err := data.Scan(&user.Username, &user.Password, &user.Privileges)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user \"%s\": %v ", username, err)
	}
	log.Println("successfully fetched user")
	return user, nil
}

func (u *Usersdb) CheckUsername(username string) (bool, error) {
	var exists bool
	data := u.db.QueryRow("SELECT exists(SELECT username from users WHERE username = ?) ", username)
	err := data.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if username \"%s\" exists: %v", username, err)
	}
	return exists, nil
}

func (u *Usersdb) GetAllUsers() ([]User, error) {
	data, err := u.db.Query("SELECT username, password, privileges from users")
	if err != nil {
		log.Fatalf("failed to fetch users: %v ", err)
	}
	defer data.Close()

	var users []User
	for data.Next() {
		var user User
		err := data.Scan(&user.Username, &user.Password, &user.Privileges)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch users: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *Usersdb) Close() error {
	err := u.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection to database: %v", err)
	}
	return nil
}
