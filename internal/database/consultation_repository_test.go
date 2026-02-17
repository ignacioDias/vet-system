package database

import (
	"testing"
	"time"
	"vetsys/internal/domain"
)

func TestConsultationRepository_CreateConsultation(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("12345678A", "John Doe", "+34600111222")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Max", "Dog", "Golden Retriever", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	consultation := domain.NewConsultation(patient.ID, "Coughing", "Kennel cough", "Antibiotics", domain.SeverityMedium)

	err := testDB.ConsultationRepo.CreateConsultation(consultation)
	if err != nil {
		t.Fatalf("Failed to create consultation: %v", err)
	}

	if consultation.ID == 0 {
		t.Error("Expected consultation ID to be set after creation")
	}
}

func TestConsultationRepository_GetConsultationByID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("23456789B", "Jane Smith", "+34600222333")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Luna", "Cat", "Persian", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	consultation := domain.NewConsultation(patient.ID, "Limping", "Sprained paw", "Rest", domain.SeverityLow)
	testDB.ConsultationRepo.CreateConsultation(consultation)

	retrieved, err := testDB.ConsultationRepo.GetConsultationByID(consultation.ID)
	if err != nil {
		t.Fatalf("Failed to get consultation by ID: %v", err)
	}

	if retrieved.Reason != consultation.Reason {
		t.Errorf("Expected reason %s, got %s", consultation.Reason, retrieved.Reason)
	}
}

func TestConsultationRepository_GetConsultationsByClientID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("34567890C", "Bob Johnson", "+34600333444")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Buddy", "Dog", "Labrador", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	cons1 := domain.NewConsultation(patient.ID, "Checkup", "Healthy", "None", domain.SeverityLow)
	testDB.ConsultationRepo.CreateConsultation(cons1)

	cons2 := domain.NewConsultation(patient.ID, "Vaccination", "Up to date", "Rabies vaccine", domain.SeverityLow)
	testDB.ConsultationRepo.CreateConsultation(cons2)

	consultations, err := testDB.ConsultationRepo.GetConsultationsByClientID(client.ID)
	if err != nil {
		t.Fatalf("Failed to get consultations by client ID: %v", err)
	}

	if len(consultations) != 2 {
		t.Errorf("Expected 2 consultations, got %d", len(consultations))
	}
}

func TestConsultationRepository_GetConsultationsByPatientID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("45678901D", "Alice Brown", "+34600444555")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Charlie", "Dog", "Beagle", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	cons := domain.NewConsultation(patient.ID, "Emergency", "Toxic ingestion", "Activated charcoal", domain.SeverityCritical)
	testDB.ConsultationRepo.CreateConsultation(cons)

	consultations, err := testDB.ConsultationRepo.GetConsultationsByPatientID(patient.ID)
	if err != nil {
		t.Fatalf("Failed to get consultations by patient ID: %v", err)
	}

	if len(consultations) != 1 {
		t.Errorf("Expected 1 consultation, got %d", len(consultations))
	}
}

func TestConsultationRepository_UpdateConsultation(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("56789012E", "Charlie Davis", "+34600666777")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2020, 4, 10, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Milo", "Dog", "Poodle", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	consultation := domain.NewConsultation(patient.ID, "Weight loss", "Under investigation", "Pending tests", domain.SeverityMedium)
	testDB.ConsultationRepo.CreateConsultation(consultation)

	consultation.Diagnosis = "Thyroid issue"
	consultation.Severity = domain.SeverityHigh
	consultation.IsCompleted = true
	consultation.UpdatedAt = time.Now()

	err := testDB.ConsultationRepo.UpdateConsultation(consultation)
	if err != nil {
		t.Fatalf("Failed to update consultation: %v", err)
	}

	retrieved, err := testDB.ConsultationRepo.GetConsultationByID(consultation.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated consultation: %v", err)
	}

	if retrieved.Diagnosis != "Thyroid issue" {
		t.Errorf("Expected updated diagnosis, got %s", retrieved.Diagnosis)
	}
}

func TestConsultationRepository_DeleteConsultation(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("67890123F", "David Wilson", "+34600777888")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Daisy", "Cat", "Maine Coon", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	consultation := domain.NewConsultation(patient.ID, "Annual checkup", "Healthy", "None", domain.SeverityLow)
	testDB.ConsultationRepo.CreateConsultation(consultation)

	err := testDB.ConsultationRepo.DeleteConsultation(consultation.ID)
	if err != nil {
		t.Fatalf("Failed to delete consultation: %v", err)
	}

	_, err = testDB.ConsultationRepo.GetConsultationByID(consultation.ID)
	if err != ErrConsultationNotFound {
		t.Errorf("Expected ErrConsultationNotFound after deletion, got %v", err)
	}
}
