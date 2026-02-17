package database

import (
	"testing"
	"vetsys/internal/domain"
)

func TestClientRepository_CreateClient(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("12345678A", "John Doe", "+34600111222")

	err := testDB.ClientRepo.CreateClient(client)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client.ID == 0 {
		t.Error("Expected client ID to be set after creation")
	}
}

func TestClientRepository_GetClientByID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("23456789B", "Jane Smith", "+34600222333")
	err := testDB.ClientRepo.CreateClient(client)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	retrieved, err := testDB.ClientRepo.GetClientByID(client.ID)
	if err != nil {
		t.Fatalf("Failed to get client by ID: %v", err)
	}

	if retrieved.ID != client.ID {
		t.Errorf("Expected client ID %d, got %d", client.ID, retrieved.ID)
	}
	if retrieved.DNI != client.DNI {
		t.Errorf("Expected DNI %s, got %s", client.DNI, retrieved.DNI)
	}
}

func TestClientRepository_GetClientByID_NotFound(t *testing.T) {
	cleanupTables(testDB)

	_, err := testDB.ClientRepo.GetClientByID(99999)
	if err != ErrClientNotFound {
		t.Errorf("Expected ErrClientNotFound, got %v", err)
	}
}

func TestClientRepository_GetClientByDNI(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("34567890C", "Bob Johnson", "+34600333444")
	err := testDB.ClientRepo.CreateClient(client)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	retrieved, err := testDB.ClientRepo.GetClientByDNI(client.DNI)
	if err != nil {
		t.Fatalf("Failed to get client by DNI: %v", err)
	}

	if retrieved.ID != client.ID {
		t.Errorf("Expected client ID %d, got %d", client.ID, retrieved.ID)
	}
}

func TestClientRepository_UpdateClient(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("45678901D", "Alice Brown", "+34600444555")
	err := testDB.ClientRepo.CreateClient(client)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	client.Name = "Alice Brown Updated"
	client.PhoneNumber = "+34600555666"
	err = testDB.ClientRepo.UpdateClient(client)
	if err != nil {
		t.Fatalf("Failed to update client: %v", err)
	}

	retrieved, err := testDB.ClientRepo.GetClientByID(client.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated client: %v", err)
	}

	if retrieved.Name != "Alice Brown Updated" {
		t.Errorf("Expected updated name, got %s", retrieved.Name)
	}
}

func TestClientRepository_DeleteClientByID(t *testing.T) {
	cleanupTables(testDB)

	client := domain.NewClient("56789012E", "Charlie Davis", "+34600666777")
	err := testDB.ClientRepo.CreateClient(client)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = testDB.ClientRepo.DeleteClientByID(client.ID)
	if err != nil {
		t.Fatalf("Failed to delete client: %v", err)
	}

	_, err = testDB.ClientRepo.GetClientByID(client.ID)
	if err != ErrClientNotFound {
		t.Errorf("Expected ErrClientNotFound after deletion, got %v", err)
	}
}

func TestClientRepository_UniqueConstraints(t *testing.T) {
	cleanupTables(testDB)

	client1 := domain.NewClient("67890123F", "David Wilson", "+34600777888")
	err := testDB.ClientRepo.CreateClient(client1)
	if err != nil {
		t.Fatalf("Failed to create first client: %v", err)
	}

	client2 := domain.NewClient("67890123F", "Different Name", "+34600888999")
	err = testDB.ClientRepo.CreateClient(client2)
	if err == nil {
		t.Error("Expected error when creating client with duplicate DNI")
	}
}
