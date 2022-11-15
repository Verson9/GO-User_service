package handler

import (
	"GO-User_service/user-service/api"
	"GO-User_service/user-service/internal/usersdb"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var _ api.ServerInterface = (*Handler)(nil)

type Database interface {
	CreateUserIfNotExists(user usersdb.User) error
	GetUser(username string) (*usersdb.User, error)
	CheckUsername(username string) (bool, error)
	GetAllUsers() ([]usersdb.User, error)
	Close() error
}

type Handler struct {
	usersdb Database
	logger  *log.Logger
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.logger.Printf("failed to send response: %v", err)
	}
}

func NewHandler(database Database, logger *log.Logger) api.ServerInterface {
	return &Handler{
		usersdb: database,
		logger:  logger,
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.usersdb.GetAllUsers()
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to fetch users from database: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(users)
	if err != nil {
		h.handleErrorResponse(w, fmt.Sprintf("failed to encode users: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, resp)
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

func (h *Handler) GetUsersUsername(w http.ResponseWriter, r *http.Request, username string) {
	userExist, err := h.usersdb.CheckUsername(username)
	if !userExist || err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, marshal)
}

func (h *Handler) handleErrorResponse(w http.ResponseWriter, msg string, errorCode int) {
	fmt.Println(msg)
	http.Error(w, msg, errorCode)
}

func writeResponse(w http.ResponseWriter, msg []byte) {
	_, err := w.Write(msg)
	if err != nil {
		fmt.Printf("failed to send response: %v", err)
		return
	}
}
