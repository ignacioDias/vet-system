package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vetsys/internal/database"
	"vetsys/internal/domain"
)

type PatientHandler struct {
	patientRepo *database.PatientRepository
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

	var patient domain.Patient
	err = json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patient.ID = idValue

	err = patientHandler.patientRepo.UpdatePatient(&patient)
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
