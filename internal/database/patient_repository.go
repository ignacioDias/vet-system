package database

import (
	"database/sql"
	"errors"
	"vetsys/internal/domain"

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
		if err == sql.ErrNoRows {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}
	return &patient, nil
}

func (patientRepository *PatientRepository) GetPatientsByOwner(ownerID int64) ([]domain.Patient, error) {
	query := `SELECT id, name, species, breed, aprox_date_of_birth, owner_id FROM patients WHERE owner_id = $1`
	var patients []domain.Patient
	err := patientRepository.DB.Select(&patients, query, ownerID)
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (patientRepository *PatientRepository) UpdatePatient(patient *domain.Patient) error {
	query := "UPDATE patients SET name = $1, species = $2, breed = $3, aprox_date_of_birth = $4, owner_id = $5 WHERE id = $6"
	result, err := patientRepository.DB.Exec(query, patient.Name, patient.Species, patient.Breed, patient.AproxDateOfBirth, patient.OwnerID, patient.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPatientNotFound
	}
	return nil
}

func (patientRepository *PatientRepository) DeletePatientByID(id int64) error {
	result, err := patientRepository.DB.Exec("DELETE FROM patients WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPatientNotFound
	}

	return nil
}
