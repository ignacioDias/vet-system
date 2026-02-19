package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vetsys/internal/database"
	"vetsys/internal/domain"
)

type ClientHandler struct {
	clientRepo *database.ClientRepository
}

func NewClientHandler(clientRepo *database.ClientRepository) *ClientHandler {
	return &ClientHandler{
		clientRepo: clientRepo,
	}
}

type ClientUpdate struct {
	DNI         *string `json:"dni"`
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phoneNumber"`
}

func (clientHandler *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client domain.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if client.DNI == "" {
		http.Error(w, "DNI is required", http.StatusBadRequest)
		return
	}
	if client.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if client.PhoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}
	err = clientHandler.clientRepo.CreateClient(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func (clientHandler *ClientHandler) GetClientByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("client_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client, err := clientHandler.clientRepo.GetClientByID(idValue)
	if err == database.ErrClientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func (clientHandler *ClientHandler) GetClientByDNIHandler(w http.ResponseWriter, r *http.Request) {
	dni := r.PathValue("client_dni")
	if dni == "" {
		http.Error(w, "No dni passed", http.StatusBadRequest)
		return
	}
	client, err := clientHandler.clientRepo.GetClientByDNI(dni)
	if err == database.ErrClientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func (clientHandler *ClientHandler) UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("client_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var clientUpdate ClientUpdate
	err = json.NewDecoder(r.Body).Decode(&clientUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client, err := clientHandler.clientRepo.GetClientByID(idValue)
	if err == database.ErrClientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if clientUpdate.DNI != nil {
		if *clientUpdate.DNI == "" {
			http.Error(w, "DNI cannot be empty", http.StatusBadRequest)
			return
		}
		client.DNI = *clientUpdate.DNI
	}
	if clientUpdate.Name != nil {
		if *clientUpdate.Name == "" {
			http.Error(w, "Name cannot be empty", http.StatusBadRequest)
			return
		}
		client.Name = *clientUpdate.Name
	}
	if clientUpdate.PhoneNumber != nil {
		if *clientUpdate.PhoneNumber == "" {
			http.Error(w, "Phone number cannot be empty", http.StatusBadRequest)
			return
		}
		client.PhoneNumber = *clientUpdate.PhoneNumber
	}
	err = clientHandler.clientRepo.UpdateClient(client)
	if err == database.ErrClientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (clientHandler *ClientHandler) DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("client_id")
	if id == "" {
		http.Error(w, "No id passed", http.StatusBadRequest)
		return
	}
	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = clientHandler.clientRepo.DeleteClientByID(idValue)
	if err == database.ErrClientNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
