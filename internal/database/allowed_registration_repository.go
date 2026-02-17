package database

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrDNIAlreadyExists = errors.New("dni already exists")
	ErrDNIInvalidOrUsed = errors.New("dni not allowed or already used")
	ErrDNINotFound      = errors.New("dni not found")
)

type AllowedRegistrationRepository struct {
	DB *sqlx.DB
}

func (r *AllowedRegistrationRepository) InsertDNI(dni string) error {
	_, err := r.DB.Exec(`INSERT INTO allowed_registrations (dni) VALUES ($1)`, dni)
	if err != nil {
		// Check for unique constraint violation (PostgreSQL error code 23505)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return ErrDNIAlreadyExists
		}
		return err
	}
	return nil
}

func (r *AllowedRegistrationRepository) UseDNI(dni string) error {
	query := `UPDATE allowed_registrations SET used = TRUE WHERE dni = $1 AND used = FALSE`
	result, err := r.DB.Exec(query, dni)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrDNIInvalidOrUsed
	}
	return nil
}

func (r *AllowedRegistrationRepository) DeleteDNI(dni string) error {
	result, err := r.DB.Exec(`DELETE FROM allowed_registrations WHERE dni = $1`, dni)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrDNINotFound
	}
	return nil
}
