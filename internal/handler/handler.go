package handler

import (
	"GO-User_service/internal/usersdb"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type database interface {
	CreateUserIfNotExists(user usersdb.User) error
	GetUser(username string) (usersdb.User, error)
	CheckUsername(username string) (bool, error)
	GetAllUsers() ([]usersdb.User, error)
	Close() error
}

type Handler struct {
	usersdb database
	logger  log.Logger
}

func NewHandler(database database) Handler {
	return Handler{
		usersdb: database,
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.usersdb.GetAllUsers()
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to fetch users from database: %v", err), http.StatusInternalServerError)
		return
	}
	respBody := new(bytes.Buffer)
	json.NewEncoder(respBody).Encode(users)
	_, err = w.Write(respBody.Bytes())
	if err != nil {
		h.logger.Println("failed to send response: %v", err)
	}
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var u usersdb.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to decode json: %v", err), http.StatusInternalServerError)
		return
	}

	err = h.usersdb.CreateUserIfNotExists(u)
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
	}
}

func (h *Handler) GetUsersByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	checkUsername, err := h.usersdb.CheckUsername(username)
	if !checkUsername || err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to get user with username [%s]: %v", username, err), http.StatusInternalServerError)
		return
	}
	user, err := h.usersdb.GetUser(username)
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to fetch user: %v", err), http.StatusInternalServerError)
		return
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to create json from user struct: %v", err), http.StatusInternalServerError)

		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to send response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleErrorResponse(w http.ResponseWriter, msg string, errorCode int) {
	fmt.Println(msg)
	http.Error(w, msg, errorCode)
}
