package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vetsys/internal/database"
	"vetsys/internal/domain"
	"vetsys/internal/utils"
)

type ConsultationHandler struct {
	consultRepo *database.ConsultationRepository
}

type ConsultationRequest struct {
	PatientID   int64           `json:"patient_id"`
	Reason      string          `json:"reason"`
	Diagnosis   string          `json:"diagnosis"`
	Treatment   string          `json:"treatment"`
	Severity    domain.Severity `json:"severity"`
	IsCompleted bool            `json:"is_completed"`
}

type ConsultationUpdate struct {
	PatientID   *int64           `json:"patient_id"`
	Reason      *string          `json:"reason"`
	Diagnosis   *string          `json:"diagnosis"`
	Treatment   *string          `json:"treatment"`
	Severity    *domain.Severity `json:"severity"`
	IsCompleted *bool            `json:"is_completed"`
}

type PaginatedResponse struct {
	Data       any   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewConsultationHandler(consultRepo *database.ConsultationRepository) *ConsultationHandler {
	return &ConsultationHandler{
		consultRepo: consultRepo,
	}
}

func (consultationHandler *ConsultationHandler) CreateConsultationHandler(w http.ResponseWriter, r *http.Request) {
	var consultationRequest ConsultationRequest
	err := json.NewDecoder(r.Body).Decode(&consultationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if consultationRequest.Reason == "" {
		http.Error(w, "Reason is required", http.StatusBadRequest)
		return
	}
	if !isValidSeverity(consultationRequest.Severity) {
		http.Error(w, "Invalid severity. Must be LOW, MEDIUM, HIGH, or CRITICAL", http.StatusBadRequest)
		return
	}

	consultation := domain.NewConsultation(consultationRequest.PatientID, consultationRequest.Reason, consultationRequest.Diagnosis, consultationRequest.Treatment, consultationRequest.Severity)
	err = consultationHandler.consultRepo.CreateConsultation(consultation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(consultation)
}

func (consultationHandler *ConsultationHandler) GetConsultationByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("consultation_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	consultation, err := consultationHandler.consultRepo.GetConsultationByID(idValue)
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
	limit, offset := utils.Pagination(r)
	total, err := consultationHandler.consultRepo.GetConsultationsByClientIDCount(idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	page := (offset / limit) + 1

	consultations, err := consultationHandler.consultRepo.GetConsultationsByClientID(idValue, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := PaginatedResponse{
		Data:       consultations,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (consultationHandler *ConsultationHandler) GetConsultationsByPatientIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("patient_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, offset := utils.Pagination(r)
	total, err := consultationHandler.consultRepo.GetConsultationsByPatientIDCount(idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	page := (offset / limit) + 1

	consultations, err := consultationHandler.consultRepo.GetConsultationsByPatientID(idValue, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := PaginatedResponse{
		Data:       consultations,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (consultationHandler *ConsultationHandler) GetAllConsultationsHandler(w http.ResponseWriter, r *http.Request) {
	isCompletedParam := r.URL.Query().Get("is_completed")

	var consultations []domain.Consultation
	var total int64
	var err error
	limit, offset := utils.Pagination(r)
	if isCompletedParam != "" {
		isCompletedValue, parseErr := strconv.ParseBool(isCompletedParam)
		if parseErr != nil {
			http.Error(w, "Invalid is_completed parameter. Use 'true' or 'false'", http.StatusBadRequest)
			return
		}
		total, err = consultationHandler.consultRepo.GetConsultationsByIsCompletedCount(isCompletedValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		consultations, err = consultationHandler.consultRepo.GetConsultationsByIsCompleted(isCompletedValue, limit, offset)
	} else {
		total, err = consultationHandler.consultRepo.GetAllConsultationsCount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		consultations, err = consultationHandler.consultRepo.GetAllConsultations(limit, offset)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	page := (offset / limit) + 1

	response := PaginatedResponse{
		Data:       consultations,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (consultationHandler *ConsultationHandler) UpdateConsultationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("consultation_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var consultationUpdate ConsultationUpdate
	err = json.NewDecoder(r.Body).Decode(&consultationUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	consultation, err := consultationHandler.consultRepo.GetConsultationByID(idValue)
	if err == database.ErrConsultationNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if consultationUpdate.PatientID != nil {
		if *consultationUpdate.PatientID <= 0 {
			http.Error(w, "PatientID must be greater than 0", http.StatusBadRequest)
			return
		}
		consultation.PatientID = *consultationUpdate.PatientID
	}
	if consultationUpdate.Reason != nil {
		if *consultationUpdate.Reason == "" {
			http.Error(w, "Reason cannot be empty", http.StatusBadRequest)
			return
		}
		consultation.Reason = *consultationUpdate.Reason
	}
	if consultationUpdate.Diagnosis != nil {
		consultation.Diagnosis = *consultationUpdate.Diagnosis
	}
	if consultationUpdate.Treatment != nil {
		consultation.Treatment = *consultationUpdate.Treatment
	}
	if consultationUpdate.Severity != nil {
		if !isValidSeverity(*consultationUpdate.Severity) {
			http.Error(w, "Invalid severity. Must be LOW, MEDIUM, HIGH, or CRITICAL", http.StatusBadRequest)
			return
		}
		consultation.Severity = *consultationUpdate.Severity
	}
	if consultationUpdate.IsCompleted != nil {
		consultation.IsCompleted = *consultationUpdate.IsCompleted
	}

	consultation.UpdatedAt = time.Now()

	err = consultationHandler.consultRepo.UpdateConsultation(consultation)
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
	id := r.PathValue("consultation_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = consultationHandler.consultRepo.DeleteConsultation(idValue)
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

func isValidSeverity(severity domain.Severity) bool {
	return severity == domain.SeverityLow ||
		severity == domain.SeverityMedium ||
		severity == domain.SeverityHigh ||
		severity == domain.SeverityCritical
}
