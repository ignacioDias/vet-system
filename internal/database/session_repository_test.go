package database

import (
	"testing"
	"time"
	"vetsys/internal/domain"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("12345678A", "test@example.com", "hashedpassword", "Test User", "profile.jpg")
	testDB.UserRepo.CreateUser(user)

	session, err := domain.NewSession(user.ID, 24*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create session object: %v", err)
	}

	err = testDB.SessionRepo.CreateSession(session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == "" {
		t.Error("Expected session ID to be set")
	}
}

func TestSessionRepository_GetSession(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("23456789B", "user1@example.com", "hashedpass1", "User One", "avatar1.jpg")
	testDB.UserRepo.CreateUser(user)

	session, _ := domain.NewSession(user.ID, 24*time.Hour)
	testDB.SessionRepo.CreateSession(session)

	retrieved, err := testDB.SessionRepo.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrieved.ID)
	}
}

func TestSessionRepository_GetSession_Expired(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("34567890C", "user2@example.com", "hashedpass2", "User Two", "avatar2.jpg")
	testDB.UserRepo.CreateUser(user)

	session, _ := domain.NewSession(user.ID, -1*time.Second)
	testDB.SessionRepo.CreateSession(session)

	_, err := testDB.SessionRepo.GetSession(session.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound for expired session, got %v", err)
	}
}

func TestSessionRepository_DeleteSessionByID(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("45678901D", "user3@example.com", "hashedpass3", "User Three", "avatar3.jpg")
	testDB.UserRepo.CreateUser(user)

	session, _ := domain.NewSession(user.ID, 24*time.Hour)
	testDB.SessionRepo.CreateSession(session)

	err := testDB.SessionRepo.DeleteSessionByID(session.ID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	_, err = testDB.SessionRepo.GetSession(session.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound after deletion, got %v", err)
	}
}

func TestSessionRepository_DeleteOldSessions(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("56789012E", "user4@example.com", "hashedpass4", "User Four", "avatar4.jpg")
	testDB.UserRepo.CreateUser(user)

	expiredSession, _ := domain.NewSession(user.ID, -1*time.Hour)
	testDB.SessionRepo.CreateSession(expiredSession)

	validSession, _ := domain.NewSession(user.ID, 24*time.Hour)
	testDB.SessionRepo.CreateSession(validSession)

	err := testDB.SessionRepo.DeleteOldSessions()
	if err != nil {
		t.Fatalf("Failed to delete old sessions: %v", err)
	}

	_, err = testDB.SessionRepo.GetSession(expiredSession.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Expected expired session to be deleted")
	}

	_, err = testDB.SessionRepo.GetSession(validSession.ID)
	if err != nil {
		t.Errorf("Expected valid session to still exist, got error: %v", err)
	}
}

func TestSessionRepository_CascadeDeleteOnUserDelete(t *testing.T) {
	cleanupTables(testDB)

	user := domain.NewUser("67890123F", "user5@example.com", "hashedpass5", "User Five", "avatar5.jpg")
	testDB.UserRepo.CreateUser(user)

	session, _ := domain.NewSession(user.ID, 24*time.Hour)
	testDB.SessionRepo.CreateSession(session)

	err := testDB.UserRepo.DeleteUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err = testDB.SessionRepo.GetSession(session.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Expected session to be cascade deleted, got error: %v", err)
	}
}
