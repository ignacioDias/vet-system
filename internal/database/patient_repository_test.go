package database

import (
	"testing"
	"time"
	"vetsys/internal/domain"
)

func TestPatientRepository_CreatePatient(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("12345678A", "John Doe", "+34600111222")
	err := testDB.ClientRepo.CreateClient(client)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	dob := time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Max", "Dog", "Golden Retriever", dob, client.ID)

	err = testDB.PatientRepo.CreatePatient(patient)
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}

	if patient.ID == 0 {
		t.Error("Expected patient ID to be set after creation")
	}
}

func TestPatientRepository_GetPatientByID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("23456789B", "Jane Smith", "+34600222333")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2019, 3, 10, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Luna", "Cat", "Persian", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	retrieved, err := testDB.PatientRepo.GetPatientByID(patient.ID)
	if err != nil {
		t.Fatalf("Failed to get patient by ID: %v", err)
	}

	if retrieved.Name != patient.Name {
		t.Errorf("Expected name %s, got %s", patient.Name, retrieved.Name)
	}
}

func TestPatientRepository_GetPatientsByOwner(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("34567890C", "Bob Johnson", "+34600333444")
	testDB.ClientRepo.CreateClient(client)

	dob1 := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	patient1 := domain.NewPatient("Buddy", "Dog", "Labrador", dob1, client.ID)
	testDB.PatientRepo.CreatePatient(patient1)

	dob2 := time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC)
	patient2 := domain.NewPatient("Charlie", "Dog", "Beagle", dob2, client.ID)
	testDB.PatientRepo.CreatePatient(patient2)

	patients, err := testDB.PatientRepo.GetPatientsByOwner(client.ID)
	if err != nil {
		t.Fatalf("Failed to get patients by owner: %v", err)
	}

	if len(patients) != 2 {
		t.Errorf("Expected 2 patients, got %d", len(patients))
	}
}

func TestPatientRepository_UpdatePatient(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("45678901D", "Alice Brown", "+34600444555")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2018, 8, 20, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Rocky", "Dog", "Bulldog", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	patient.Name = "Rocky Jr."
	patient.Breed = "French Bulldog"
	err := testDB.PatientRepo.UpdatePatient(patient)
	if err != nil {
		t.Fatalf("Failed to update patient: %v", err)
	}

	retrieved, err := testDB.PatientRepo.GetPatientByID(patient.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated patient: %v", err)
	}

	if retrieved.Name != "Rocky Jr." {
		t.Errorf("Expected updated name, got %s", retrieved.Name)
	}
}

func TestPatientRepository_DeletePatientByID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("56789012E", "Charlie Davis", "+34600666777")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2019, 12, 5, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Bella", "Cat", "Siamese", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	err := testDB.PatientRepo.DeletePatientByID(patient.ID)
	if err != nil {
		t.Fatalf("Failed to delete patient: %v", err)
	}

	_, err = testDB.PatientRepo.GetPatientByID(patient.ID)
	if err != ErrPatientNotFound {
		t.Errorf("Expected ErrPatientNotFound after deletion, got %v", err)
	}
}

func TestPatientRepository_CascadeDelete(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("67890123F", "David Wilson", "+34600777888")
	testDB.ClientRepo.CreateClient(client)

	dob := time.Date(2020, 4, 10, 0, 0, 0, 0, time.UTC)
	patient := domain.NewPatient("Milo", "Dog", "Poodle", dob, client.ID)
	testDB.PatientRepo.CreatePatient(patient)

	err := testDB.ClientRepo.DeleteClientByID(client.ID)
	if err != nil {
		t.Fatalf("Failed to delete client: %v", err)
	}

	_, err = testDB.PatientRepo.GetPatientByID(patient.ID)
	if err != ErrPatientNotFound {
		t.Errorf("Expected patient to be cascade deleted, got error: %v", err)
	}
}
