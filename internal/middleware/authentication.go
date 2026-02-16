package middleware

import (
	"context"
	"net/http"
	"vetsys/internal/database"
)

type AuthMiddleware struct {
	SessionRepo *database.SessionRepository
}

type contextKey string

const UserIDKey contextKey = "userID"

func (auth *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		session, err := auth.SessionRepo.GetSession(cookie.Value)
		if err == database.ErrSessionNotFound {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
