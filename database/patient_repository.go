package database

import (
	"errors"
	"vetsys/domain"

	"github.com/jmoiron/sqlx"
)

type PatientRepository struct {
	DB *sqlx.DB
}

var ErrPatientNotFound = errors.New("Patient not found")

func (patientRepository *PatientRepository) CreatePatient(patient *domain.Patient) error {
	query := `
	INSERT INTO patients (name, species, breed, aprox_date_of_birth, owner_id)
	VALUES (:name, :species, :breed, :aprox_date_of_birth, :owner_id)
	RETURNING id
	`
	stmt, err := patientRepository.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(&patient.ID, patient)
}

func (patientRepository *PatientRepository) GetPatientByID(id int64) (*domain.Patient, error) {
	query := `SELECT id, name, species, breed, aprox_date_of_birth, owner_id FROM patients WHERE id = $1`
	patient := domain.Patient{}
	err := patientRepository.DB.Get(&patient, query, id)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (patientRepository *PatientRepository) GetPatientsByOwner(ownerID int64) (*[]domain.Patient, error) {
	query := `SELECT id, name, species, breed, aprox_date_of_birth, owner_id FROM patients WHERE owner_id = $1`
	rows, err := patientRepository.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patients []domain.Patient
	for rows.Next() {
		var patient domain.Patient
		if err := rows.Scan(&patient.ID, &patient.Name, &patient.Breed, &patient.AproxDateOfBirth, &patient.OwnerID); err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return &patients, nil
}
