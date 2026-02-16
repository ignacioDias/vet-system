package domain

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Session struct {
	ID        string    `json:"id" db:"id"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UserID    int64     `json:"userId" db:"user_id"`
}

func NewSession(userID int64, duration time.Duration) (*Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &Session{
		ID:        id,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
	}, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32) // 256 bits of entropy

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
