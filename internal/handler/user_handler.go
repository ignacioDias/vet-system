package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vetsys/internal/database"
	"vetsys/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepo *database.UserRepository
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

func (UserHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (userHandler *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidPassword(user.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	err = userHandler.UserRepo.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (userHandler *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = userHandler.UserRepo.DeleteUserByID(idValue)
	if err == database.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (userHandler *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user.ID = idValue
	err = userHandler.UserRepo.UpdateUser(&user)
	if err == database.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (userHandler *UserHandler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !isValidPassword(req.Password) {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	err = userHandler.UserRepo.UpdatePassword(idValue, string(hashedPassword))
	if err == database.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
