package handler

import (
	"net/http"
	"vetsys/internal/database"
)

type ClientHandler struct {
	clientRepo *database.ClientRepository
}

type UserResponse struct {
	ID             int64  `db:"id" json:"id"`
	DNI            string `db:"dni" json:"dni"`
	Email          string `db:"email" json:"email"`
	Name           string `db:"name" json:"name"`
	ProfilePicture string `db:"profile_picture" json:"profilePicture"`
}

func (clientHandler *ClientHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}
