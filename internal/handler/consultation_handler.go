package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vetsys/internal/database"
	"vetsys/internal/domain"
)

type ConsultationHandler struct {
	consultationRepo *database.ConsultationRepository
}

type ConsultationRequest struct {
	PatientID   int64           `db:"patient_id" json:"patient_id"`
	Reason      string          `db:"reason" json:"reason"`
	Diagnosis   string          `db:"diagnosis" json:"diagnosis"`
	Treatment   string          `db:"treatment" json:"treatment"`
	Severity    domain.Severity `db:"severity" json:"severity"`
	IsCompleted bool            `db:"is_completed" json:"is_completed"`
}

func (consultationHandler *ConsultationHandler) CreateConsultationHandler(w http.ResponseWriter, r *http.Request) {
	var consultationRequest ConsultationRequest
	err := json.NewDecoder(r.Body).Decode(&consultationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	consultation := domain.NewConsultation(consultationRequest.PatientID, consultationRequest.Reason, consultationRequest.Diagnosis, consultationRequest.Treatment, consultationRequest.Severity)
	err = consultationHandler.consultationRepo.CreateConsultation(consultation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(consultation)
}

func (consultationHandler *ConsultationHandler) GetConsultationByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("consultation-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	consultation, err := consultationHandler.consultationRepo.GetConsultationByID(idValue)
	if err == database.ErrConsultationNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultation)
}

func (consultationHandler *ConsultationHandler) GetConsultationsByClientIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("client-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	consultations, err := consultationHandler.consultationRepo.GetConsultationsByClientID(idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultations)
}

func (consultationHandler *ConsultationHandler) GetConsultationsByPatientIDHandler(w http.ResponseWriter, r *http.Request) {
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
	consultations, err := consultationHandler.consultationRepo.GetConsultationsByPatientID(idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultations)
}

func (consultationHandler *ConsultationHandler) GetAllConsultationsHandler(w http.ResponseWriter, r *http.Request) {
	isCompletedParam := r.URL.Query().Get("is_completed")

	var consultations []domain.Consultation
	var err error

	if isCompletedParam != "" {
		isCompletedValue, parseErr := strconv.ParseBool(isCompletedParam)
		if parseErr != nil {
			http.Error(w, "Invalid is_completed parameter. Use 'true' or 'false'", http.StatusBadRequest)
			return
		}
		consultations, err = consultationHandler.consultationRepo.GetConsultationsByIsCompleted(isCompletedValue)
	} else {
		consultations, err = consultationHandler.consultationRepo.GetAllConsultations()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultations)
}
func (consultationHandler *ConsultationHandler) UpdateConsultationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("consultation-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var consultationRequest ConsultationRequest
	err = json.NewDecoder(r.Body).Decode(&consultationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	consultation := domain.Consultation{
		ID:          idValue,
		PatientID:   consultationRequest.PatientID,
		Reason:      consultationRequest.Reason,
		Diagnosis:   consultationRequest.Diagnosis,
		Treatment:   consultationRequest.Treatment,
		Severity:    consultationRequest.Severity,
		IsCompleted: consultationRequest.IsCompleted,
		UpdatedAt:   time.Now(),
	}
	err = consultationHandler.consultationRepo.UpdateConsultation(&consultation)
	if err == database.ErrConsultationNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (consultationHandler *ConsultationHandler) DeleteConsultationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("consultation-id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = consultationHandler.consultationRepo.DeleteConsultation(idValue)
	if err == database.ErrConsultationNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
