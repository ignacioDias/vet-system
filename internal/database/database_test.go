package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var testDB *DataBase

func TestMain(m *testing.M) {
	// Setup test database connection
	db, err := setupTestDB()
	if err != nil {
		fmt.Printf("Failed to setup test database: %v\n", err)
		os.Exit(1)
	}
	testDB = db

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestDB(testDB)

	os.Exit(code)
}

func setupTestDB() (*DataBase, error) {
	// Use environment variable or default test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/vetsys_test?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Create database instance
	database := NewDataBase(db)

	// Initialize tables
	if err := database.Init(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return database, nil
}

func cleanupTestDB(db *DataBase) {
	if db != nil && db.DB != nil {
		// Drop all tables in reverse order of dependencies
		db.DB.Exec("DROP TABLE IF EXISTS sessions CASCADE")
		db.DB.Exec("DROP TABLE IF EXISTS consultations CASCADE")
		db.DB.Exec("DROP TABLE IF EXISTS patients CASCADE")
		db.DB.Exec("DROP TABLE IF EXISTS clients CASCADE")
		db.DB.Exec("DROP TABLE IF EXISTS users CASCADE")
		db.DB.Exec("DROP TABLE IF EXISTS allowed_registrations CASCADE")
		db.DB.Close()
	}
}

func cleanupTables(db *DataBase) {
	// Clean tables but preserve schema
	db.DB.Exec("TRUNCATE TABLE sessions CASCADE")
	db.DB.Exec("TRUNCATE TABLE consultations CASCADE")
	db.DB.Exec("TRUNCATE TABLE patients CASCADE")
	db.DB.Exec("TRUNCATE TABLE clients CASCADE")
	db.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	db.DB.Exec("TRUNCATE TABLE allowed_registrations CASCADE")
}

func TestNewDataBase(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	if testDB.DB == nil {
		t.Error("Database connection is nil")
	}

	if testDB.UserRepo == nil {
		t.Error("UserRepository is nil")
	}

	if testDB.ClientRepo == nil {
		t.Error("ClientRepository is nil")
	}

	if testDB.PatientRepo == nil {
		t.Error("PatientRepository is nil")
	}

	if testDB.ConsultationRepo == nil {
		t.Error("ConsultationRepository is nil")
	}

	if testDB.SessionRepo == nil {
		t.Error("SessionRepository is nil")
	}

	if testDB.AllowedRegistrationsRepo == nil {
		t.Error("AllowedRegistrationsRepo is nil")
	}
}

func TestDataBaseInit(t *testing.T) {
	// Test that tables exist
	var tableNames []string
	expectedTables := []string{"users", "clients", "patients", "consultations", "sessions", "allowed_registrations"}

	query := `
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename IN ('users', 'clients', 'patients', 'consultations', 'sessions', 'allowed_registrations')
	`

	err := testDB.DB.Select(&tableNames, query)
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}

	if len(tableNames) != len(expectedTables) {
		t.Errorf("Expected %d tables, got %d", len(expectedTables), len(tableNames))
	}

	// Verify each expected table exists
	tableMap := make(map[string]bool)
	for _, table := range tableNames {
		tableMap[table] = true
	}

	for _, expected := range expectedTables {
		if !tableMap[expected] {
			t.Errorf("Expected table %s not found", expected)
		}
	}
}
