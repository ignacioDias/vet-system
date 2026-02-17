package database

import (
	"testing"
)

func TestAllowedRegistrationRepository_InsertDNI(t *testing.T) {
	cleanupTables(testDB)

	dni := "12345678A"
	err := testDB.AllowedRegistrationsRepo.InsertDNI(dni)
	if err != nil {
		t.Fatalf("Failed to insert DNI: %v", err)
	}
}

func TestAllowedRegistrationRepository_InsertDNI_Duplicate(t *testing.T) {
	cleanupTables(testDB)

	dni := "23456789B"
	err := testDB.AllowedRegistrationsRepo.InsertDNI(dni)
	if err != nil {
		t.Fatalf("Failed to insert DNI first time: %v", err)
	}

	err = testDB.AllowedRegistrationsRepo.InsertDNI(dni)
	if err != ErrDNIAlreadyExists {
		t.Errorf("Expected ErrDNIAlreadyExists, got %v", err)
	}
}

func TestAllowedRegistrationRepository_UseDNI(t *testing.T) {
	cleanupTables(testDB)

	dni := "34567890C"
	err := testDB.AllowedRegistrationsRepo.InsertDNI(dni)
	if err != nil {
		t.Fatalf("Failed to insert DNI: %v", err)
	}

	err = testDB.AllowedRegistrationsRepo.UseDNI(dni)
	if err != nil {
		t.Fatalf("Failed to use DNI: %v", err)
	}
}

func TestAllowedRegistrationRepository_UseDNI_AlreadyUsed(t *testing.T) {
	cleanupTables(testDB)

	dni := "45678901D"
	testDB.AllowedRegistrationsRepo.InsertDNI(dni)
	testDB.AllowedRegistrationsRepo.UseDNI(dni)

	err := testDB.AllowedRegistrationsRepo.UseDNI(dni)
	if err != ErrDNIInvalidOrUsed {
		t.Errorf("Expected ErrDNIInvalidOrUsed, got %v", err)
	}
}

func TestAllowedRegistrationRepository_UseDNI_NotFound(t *testing.T) {
	cleanupTables(testDB)

	err := testDB.AllowedRegistrationsRepo.UseDNI("NONEXISTENT")
	if err != ErrDNIInvalidOrUsed {
		t.Errorf("Expected ErrDNIInvalidOrUsed, got %v", err)
	}
}

func TestAllowedRegistrationRepository_DeleteDNI(t *testing.T) {
	cleanupTables(testDB)

	dni := "56789012E"
	testDB.AllowedRegistrationsRepo.InsertDNI(dni)

	err := testDB.AllowedRegistrationsRepo.DeleteDNI(dni)
	if err != nil {
		t.Fatalf("Failed to delete DNI: %v", err)
	}
}

func TestAllowedRegistrationRepository_DeleteDNI_NotFound(t *testing.T) {
	cleanupTables(testDB)

	err := testDB.AllowedRegistrationsRepo.DeleteDNI("NONEXISTENT")
	if err != ErrDNINotFound {
		t.Errorf("Expected ErrDNINotFound, got %v", err)
	}
}

func TestAllowedRegistrationRepository_Workflow(t *testing.T) {
	cleanupTables(testDB)

	dni := "67890123F"

	// Insert DNI
	err := testDB.AllowedRegistrationsRepo.InsertDNI(dni)
	if err != nil {
		t.Fatalf("Step 1 failed - Insert DNI: %v", err)
	}

	// Use the DNI
	err = testDB.AllowedRegistrationsRepo.UseDNI(dni)
	if err != nil {
		t.Fatalf("Step 2 failed - Use DNI: %v", err)
	}

	// Try to use again (should fail)
	err = testDB.AllowedRegistrationsRepo.UseDNI(dni)
	if err != ErrDNIInvalidOrUsed {
		t.Error("Step 3 failed - DNI should not be allowed after use")
	}
}
