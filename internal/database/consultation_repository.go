package database

import (
	"errors"
	"vetsys/domain"

	"github.com/jmoiron/sqlx"
)

type ConsultationRepository struct {
	DB *sqlx.DB
}

var ErrConsultationNotFound = errors.New("Consultation not found")

func (consultationRepository *ConsultationRepository) CreateConsultation(consultation *domain.Consultation) error {
	query := `INSERT INTO consultations (patient_id, reason, diagnosis, treatment, severity, is_completed, created_at, updated_at) 
	VALUES (:patient_id, :reason, :diagnosis, :treatment, :severity, :is_completed, :created_at, :updated_at) 
	RETURNING id`
	stmt, err := consultationRepository.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.Get(&consultation.ID, consultation)
}
func (consultationRepository *ConsultationRepository) GetConsultationByID(id int64) (*domain.Consultation, error) {
	consultation := domain.Consultation{}
	query := `SELECT id, patient_id, reason, diagnosis, treatment, severity, is_completed, created_at, updated_at FROM consultations WHERE id = $1`
	err := consultationRepository.DB.Get(&consultation, query, id)
	if err != nil {
		return nil, err
	}
	return &consultation, nil
}
func (consultationRepository *ConsultationRepository) GetConsultationsByClient(clientID int64) ([]domain.Consultation, error) {
	query := `
	SELECT c.id, c.patient_id, c.reason, c.diagnosis, 
	       c.treatment, c.severity, c.is_completed, 
	       c.created_at, c.updated_at
	FROM consultations c
	JOIN patients p ON c.patient_id = p.id
	WHERE p.owner_id = $1
	ORDER BY c.created_at DESC
	`
	var consultations []domain.Consultation
	err := consultationRepository.DB.Select(&consultations, query, clientID)
	if err != nil {
		return nil, err
	}
	return consultations, nil
}
func (consultationRepository *ConsultationRepository) GetConsultationsByPatient(patientID int64) ([]domain.Consultation, error) {
	query := `SELECT id, patient_id, reason, diagnosis, treatment, severity, is_completed, created_at, updated_at FROM consultations WHERE patient_id = $1`
	var consultations []domain.Consultation
	err := consultationRepository.DB.Select(&consultations, query, patientID)
	if err != nil {
		return nil, err
	}
	return consultations, nil
}
func (consultationRepository *ConsultationRepository) GetAllConsultations() ([]domain.Consultation, error) {
	query := `SELECT id, patient_id, reason, diagnosis, treatment, severity, is_completed, created_at, updated_at FROM consultations`
	var consultations []domain.Consultation
	err := consultationRepository.DB.Select(&consultations, query)
	if err != nil {
		return nil, err
	}
	return consultations, nil
}
func (consultationRepository *ConsultationRepository) GetConsultationsByIsCompleted(isCompleted bool) ([]domain.Consultation, error) {
	query := `SELECT id, patient_id, reason, diagnosis, treatment, severity, is_completed, created_at, updated_at FROM consultations WHERE is_completed = $1`
	var consultations []domain.Consultation
	err := consultationRepository.DB.Select(&consultations, query, isCompleted)
	if err != nil {
		return nil, err
	}
	return consultations, nil
}
func (consultationRepository *ConsultationRepository) UpdateConsultation(consultation *domain.Consultation) error {
	query := "UPDATE consultations SET patient_id = $1, reason = $2, diagnosis = $3, treatment = $4, severity = $5, is_completed = $6, updated_at = $7 WHERE id = $8"
	result, err := consultationRepository.DB.Exec(query, consultation.PatientID, consultation.Reason, consultation.Diagnosis, consultation.Treatment, consultation.Severity, consultation.IsCompleted, consultation.UpdatedAt, consultation.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrConsultationNotFound
	}
	return nil
}
func (consultationRepository *ConsultationRepository) DeleteConsultation(id int64) error {
	query := `DELETE FROM consultations WHERE id = $1`
	result, err := consultationRepository.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrConsultationNotFound
	}
	return nil
}
