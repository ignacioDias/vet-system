package domain

import "time"

type Consultation struct {
	ID          int64
	AnimalID    int64
	Reason      string
	Diagnosis   string
	Treatment   string
	Severity    Severity
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

func NewConsultation(animalID int64, reason string, diagnosis string, treatment string, severity Severity) *Consultation {
	return &Consultation{
		AnimalID:    animalID,
		Reason:      reason,
		Diagnosis:   diagnosis,
		Treatment:   treatment,
		Severity:    severity,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
