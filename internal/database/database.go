package database

import "github.com/jmoiron/sqlx"

type DataBase struct {
	DB               *sqlx.DB
	UserRepo         *UserRepository
	PatientRepo      *PatientRepository
	ClientRepo       *ClientRepository
	ConsultationRepo *ConsultationRepository
	SessionRepo      *SessionRepository
}

var createUserTable string = `
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    dni TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    profile_picture TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_dni ON users(dni);
`

var createClientsTable string = `
CREATE TABLE IF NOT EXISTS clients (
    id BIGSERIAL PRIMARY KEY,
    dni TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    phone_number TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_clients_dni ON clients(dni);
`

var createPatientsTable string = `
CREATE TABLE IF NOT EXISTS patients (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    breed TEXT NOT NULL,
    aprox_date_of_birth TIMESTAMP NOT NULL,
    owner_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_patients_owner_id ON patients(owner_id);
`

var createConsultationsTable string = `
CREATE TABLE IF NOT EXISTS consultations (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    diagnosis TEXT NOT NULL,
    treatment TEXT NOT NULL,
    severity TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_consultations_patient_id ON consultations(patient_id);
`

var createSessionsTable string = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
`

func NewDataBase(db *sqlx.DB) *DataBase {
	return &DataBase{
		DB:               db,
		UserRepo:         &UserRepository{DB: db},
		PatientRepo:      &PatientRepository{DB: db},
		ClientRepo:       &ClientRepository{DB: db},
		ConsultationRepo: &ConsultationRepository{DB: db},
		SessionRepo:      &SessionRepository{DB: db},
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

	_, err = d.DB.Exec(createSessionsTable)
	if err != nil {
		return err
	}
	return nil
}
