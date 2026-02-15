package database

import "github.com/jmoiron/sqlx"

type DataBase struct {
	DB               *sqlx.DB
	UserRepo         *UserRepository
	PatientRepo      *PatientRepository
	ClientRepo       *ClientRepository
	ConsultationRepo *ConsultationRepository
}

var createUserTable string = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dni TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    profilePicture TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_dni ON users(dni);
`

func NewDataBase(db *sqlx.DB) *DataBase {
	return &DataBase{
		DB:               db,
		UserRepo:         &UserRepository{DB: db},
		PatientRepo:      &PatientRepository{DB: db},
		ClientRepo:       &ClientRepository{DB: db},
		ConsultationRepo: &ConsultationRepository{DB: db},
	}
}

func (d *DataBase) Init() error {

	return nil
}
