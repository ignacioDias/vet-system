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

var createClientsTable string = `
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dni TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    phone_number TEXT NOT NULL
);
CREATE INDEX idx_clients_dni ON clients(dni);
`

var createPatientsTable string = `
CREATE TABLE IF NOT EXISTS patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    breed TEXT NOT NULL,
    aprox_date_of_birth DATETIME NOT NULL,
    owner_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES clients(id) ON DELETE CASCADE
);
CREATE INDEX idx_patients_owner_id ON patients(owner_id);
`

var createConsultationsTable string = `
CREATE TABLE IF NOT EXISTS consultations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    diagnosis TEXT NOT NULL,
    treatment TEXT NOT NULL,
    severity TEXT NOT NULL,
    is_completed INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
CREATE INDEX idx_consultations_patient_id ON consultations(patient_id);
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
	_, err := d.DB.Exec(createUserTable)
	if err != nil {
		return err
	}

	_, err = d.DB.Exec(createClientsTable)
	if err != nil {
		return err
	}

	_, err = d.DB.Exec(createPatientsTable)
	if err != nil {
		return err
	}

	_, err = d.DB.Exec(createConsultationsTable)
	if err != nil {
		return err
	}

	return nil
}
