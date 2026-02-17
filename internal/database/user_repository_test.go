package database

import (
	"testing"
	"vetsys/internal/domain"
)

func TestUserRepository_CreateUser(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("12345678A", "test@example.com", "hashedpassword", "Test User", "profile.jpg")

	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set after creation")
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("23456789B", "user1@example.com", "hashedpass1", "User One", "avatar1.jpg")
	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := testDB.UserRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, retrieved.ID)
	}
	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestUserRepository_GetUserByDNI(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("34567890C", "user2@example.com", "hashedpass2", "User Two", "avatar2.jpg")
	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := testDB.UserRepo.GetUserByDNI(user.DNI)
	if err != nil {
		t.Fatalf("Failed to get user by DNI: %v", err)
	}

	if retrieved.DNI != user.DNI {
		t.Errorf("Expected DNI %s, got %s", user.DNI, retrieved.DNI)
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("45678901D", "user3@example.com", "hashedpass3", "User Three", "avatar3.jpg")
	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrieved, err := testDB.UserRepo.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("56789012E", "user4@example.com", "hashedpass4", "User Four", "avatar4.jpg")
	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user.Name = "User Four Updated"
	user.Email = "updated4@example.com"

	err = testDB.UserRepo.UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	retrieved, err := testDB.UserRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}

	if retrieved.Name != "User Four Updated" {
		t.Errorf("Expected updated name, got %s", retrieved.Name)
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("67890123F", "user5@example.com", "oldpassword", "User Five", "avatar5.jpg")
	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	newPassword := "newhashedpassword"
	err = testDB.UserRepo.UpdatePassword(user.ID, newPassword)
	if err != nil {
		t.Fatalf("Failed to update password: %v", err)
	}

	retrieved, err := testDB.UserRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrieved.Password != newPassword {
		t.Errorf("Expected password %s, got %s", newPassword, retrieved.Password)
	}
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("78901234G", "user6@example.com", "hashedpass6", "User Six", "avatar6.jpg")
	err := testDB.UserRepo.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = testDB.UserRepo.DeleteUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err = testDB.UserRepo.GetUserByID(user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}

func TestUserRepository_UniqueConstraints(t *testing.T) {
	cleanupTables(testDB)

	user1 := domain.NewUser("89012345H", "unique@example.com", "pass1", "User Seven", "avatar7.jpg")
	err := testDB.UserRepo.CreateUser(user1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Duplicate DNI
	user2 := domain.NewUser("89012345H", "different@example.com", "pass2", "User Eight", "avatar8.jpg")
	err = testDB.UserRepo.CreateUser(user2)
	if err == nil {
		t.Error("Expected error when creating user with duplicate DNI")
	}

	// Duplicate email
	user3 := domain.NewUser("90123456I", "unique@example.com", "pass3", "User Nine", "avatar9.jpg")
	err = testDB.UserRepo.CreateUser(user3)
	if err == nil {
		t.Error("Expected error when creating user with duplicate email")
	}
}
