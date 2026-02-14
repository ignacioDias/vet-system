package database

import "github.com/jmoiron/sqlx"

type ConsultationRepository struct {
	DB *sqlx.DB
}
