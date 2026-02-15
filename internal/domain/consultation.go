package domain

import "time"

type Consultation struct {
	ID          int64     `db:"id" json:"id"`
	PatientID   int64     `db:"patient_id" json:"patient_id"`
	Reason      string    `db:"reason" json:"reason"`
	Diagnosis   string    `db:"diagnosis" json:"diagnosis"`
	Treatment   string    `db:"treatment" json:"treatment"`
	Severity    Severity  `db:"severity" json:"severity"`
	IsCompleted bool      `db:"is_completed" json:"is_completed"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

func NewConsultation(patientID int64, reason string, diagnosis string, treatment string, severity Severity) *Consultation {
	return &Consultation{
		PatientID:   patientID,
		Reason:      reason,
		Diagnosis:   diagnosis,
		Treatment:   treatment,
		Severity:    severity,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
