package main

import (
	"log"
	"os"
	"vetsys/internal/database"
	"vetsys/internal/handler"
	"vetsys/internal/router"
	"vetsys/internal/server"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	sqlxDB, err := sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlxDB.Close()

	db := database.NewDataBase(sqlxDB)
	err = db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	clientHandler := handler.NewClientHandler(db.ClientRepo)
	consultHandler := handler.NewConsultationHandler(db.ConsultationRepo)
	patientHandler := handler.NewPatientHandler(db.PatientRepo)
	userHandler := handler.NewUserHandler(db.UserRepo, db.SessionRepo)

	r := router.NewRouter(clientHandler, consultHandler, patientHandler, userHandler)
	srv := server.NewServer("8888", r)
	srv.StartServer(*r)
}
