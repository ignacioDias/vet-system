package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vetsys/internal/database"
	"vetsys/internal/domain"
)

type PatientHandler struct {
	patientRepo *database.PatientRepository
}

type PatientUpdate struct {
	Name             *string    `json:"name"`
	Species          *string    `json:"species"`
	Breed            *string    `json:"breed"`
	AproxDateOfBirth *time.Time `json:"aproxDateOfBirth"`
	OwnerID          *int64     `json:"ownerId"`
}

func (patientHandler *PatientHandler) CreatePatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient domain.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = patientHandler.patientRepo.CreatePatient(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

func (patientHandler *PatientHandler) GetPatientByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("patient-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patient, err := patientHandler.patientRepo.GetPatientByID(idValue)
	if err == database.ErrPatientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(patient)
}
func (patientHandler *PatientHandler) GetPatientByOwnerIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("owner-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patients, err := patientHandler.patientRepo.GetPatientsByOwner(idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(patients)
}

func (patientHandler *PatientHandler) UpdatePatientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("patient-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var patientUpdate PatientUpdate
	err = json.NewDecoder(r.Body).Decode(&patientUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patient, err := patientHandler.patientRepo.GetPatientByID(idValue)
	if err == database.ErrPatientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if patientUpdate.Name != nil {
		if *patientUpdate.Name == "" {
			http.Error(w, "Name cannot be empty", http.StatusBadRequest)
			return
		}
		patient.Name = *patientUpdate.Name
	}
	if patientUpdate.Species != nil {
		if *patientUpdate.Species == "" {
			http.Error(w, "Species cannot be empty", http.StatusBadRequest)
			return
		}
		patient.Species = *patientUpdate.Species
	}
	if patientUpdate.Breed != nil {
		if *patientUpdate.Breed == "" {
			http.Error(w, "Breed cannot be empty", http.StatusBadRequest)
			return
		}
		patient.Breed = *patientUpdate.Breed
	}
	if patientUpdate.AproxDateOfBirth != nil {
		patient.AproxDateOfBirth = *patientUpdate.AproxDateOfBirth
	}
	if patientUpdate.OwnerID != nil {
		if *patientUpdate.OwnerID <= 0 {
			http.Error(w, "OwnerID must be greater than 0", http.StatusBadRequest)
			return
		}
		patient.OwnerID = *patientUpdate.OwnerID
	}

	err = patientHandler.patientRepo.UpdatePatient(patient)
	if err == database.ErrPatientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (patientHandler *PatientHandler) DeletePatientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("patient-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = patientHandler.patientRepo.DeletePatientByID(idValue)
	if err == database.ErrPatientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
