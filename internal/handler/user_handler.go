package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
	"vetsys/internal/database"
	"vetsys/internal/domain"
	"vetsys/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepo          *database.UserRepository
	SessionRepo       *database.SessionRepository
	AllowedRegistRepo *database.AllowedRegistrationRepository
}

func NewUserHandler(userRepo *database.UserRepository, sessionRepo *database.SessionRepository, allowedRegistrationsRepo *database.AllowedRegistrationRepository) *UserHandler {
	return &UserHandler{
		UserRepo:          userRepo,
		SessionRepo:       sessionRepo,
		AllowedRegistRepo: allowedRegistrationsRepo,
	}
}

var isProduction bool = os.Getenv("ENV") == "production"

type CreateUserRequest struct {
	DNI            string `json:"dni"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
}
type LoginRequest struct {
	DNI      string `json:"dni"`
	Password string `json:"password"`
}
type UpdatePasswordRequest struct {
	Password string `json:"password"`
}
type UserUpdate struct {
	Email          *string `json:"email"`
	Name           *string `json:"name"`
	ProfilePicture *string `json:"profilePicture"`
}

func (UserHandler *UserHandler) LogInHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := UserHandler.UserRepo.GetUserByDNI(loginRequest.DNI)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	session, err := domain.NewSession(user.ID, 24*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = UserHandler.SessionRepo.CreateSession(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction, // false only in localhost dev
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (userHandler *UserHandler) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	err = userHandler.SessionRepo.DeleteSessionByID(cookie.Value)
	if err != nil && err != database.ErrSessionNotFound {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"logged out"}`))
}

func (userHandler *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if req.DNI == "" {
		http.Error(w, "DNI is required", http.StatusBadRequest)
		return
	}

	if req.ProfilePicture == "" {
		req.ProfilePicture = "https://oyster.ignimgs.com/mediawiki/apis.ign.com/adventure-time-hey-ice-king/a/a6/JakeHeadshot.jpg"
	}

	if !isValidPassword(req.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	err = userHandler.AllowedRegistRepo.UseDNI(req.DNI)
	if err != nil {
		http.Error(w, "Invalid DNI", http.StatusBadRequest)
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
	req.Password = string(hashedPassword)

	user := domain.NewUser(req.DNI, req.Email, req.Password, req.Name, req.ProfilePicture)

	err = userHandler.UserRepo.CreateUser(user)
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
	idValue, ok := userHandler.authorizeUserAccess(w, r)
	if !ok {
		return
	}
	err := userHandler.UserRepo.DeleteUserByID(idValue)
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
	idValue, ok := userHandler.authorizeUserAccess(w, r)
	if !ok {
		return
	}

	var userUpdate UserUpdate
	err := json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := userHandler.UserRepo.GetUserByID(idValue)
	if err == database.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userUpdate.Email != nil {
		if *userUpdate.Email == "" {
			http.Error(w, "Email cannot be empty", http.StatusBadRequest)
			return
		}
		if !isValidEmail(*userUpdate.Email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}
		user.Email = *userUpdate.Email
	}
	if userUpdate.Name != nil {
		if *userUpdate.Name == "" {
			http.Error(w, "Name cannot be empty", http.StatusBadRequest)
			return
		}
		user.Name = *userUpdate.Name
	}
	if userUpdate.ProfilePicture != nil {
		user.ProfilePicture = *userUpdate.ProfilePicture
	}

	err = userHandler.UserRepo.UpdateUser(user)
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
	idValue, ok := userHandler.authorizeUserAccess(w, r)
	if !ok {
		return
	}

	var req UpdatePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
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
	userHandler.LogOutHandler(w, r)
}

func (userHandler *UserHandler) authorizeUserAccess(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id := r.PathValue("user_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return 0, false
	}

	pathUserID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, false
	}

	sessionUserID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0, false
	}

	if sessionUserID != pathUserID {
		http.Error(w, "Forbidden: You can only modify your own account", http.StatusForbidden)
		return 0, false
	}
	return pathUserID, true
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
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
