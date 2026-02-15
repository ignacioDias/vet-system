package handler

import "vetsys/internal/database"

type UserHandler struct {
	UserRepo *database.UserRepository
}
type UserResponse struct {
	ID             int64  `db:"id" json:"id"`
	DNI            string `db:"dni" json:"dni"`
	Email          string `db:"email" json:"email"`
	Name           string `db:"name" json:"name"`
	ProfilePicture string `db:"profile_picture" json:"profilePicture"`
}
