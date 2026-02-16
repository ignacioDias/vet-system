package database

import (
	"database/sql"
	"errors"
	"vetsys/internal/domain"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	DB *sqlx.DB
}

var ErrSessionNotFound = errors.New("Session not found")

func (sessionRepo *SessionRepository) CreateSession(session *domain.Session) error {
	query := `
	INSERT INTO sessions (id, expires_at, created_at, user_id)
	VALUES (:id, :expires_at, :created_at, :user_id)
	`
	_, err := sessionRepo.DB.NamedExec(query, session)
	return err
}

func (sessionRepo *SessionRepository) GetSession(id string) (*domain.Session, error) {
	query := `SELECT id, expires_at, created_at, user_id FROM sessions WHERE id = $1 AND expires_at > (now() AT TIME ZONE 'UTC')`
	session := domain.Session{}
	err := sessionRepo.DB.Get(&session, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}
func (sessionRepo *SessionRepository) DeleteSessionByID(id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	result, err := sessionRepo.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (sessionRepo *SessionRepository) DeleteOldSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < (now() AT TIME ZONE 'UTC')`
	result, err := sessionRepo.DB.Exec(query)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
