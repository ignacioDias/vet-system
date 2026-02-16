package router

import (
	"net/http"
	"vetsys/internal/handler"
	"vetsys/internal/middleware"
)

type Router struct {
	mux                 *http.ServeMux
	clientHandler       *handler.ClientHandler
	consultationHandler *handler.ConsultationHandler
	patientHandler      *handler.PatientHandler
	userHandler         *handler.UserHandler
	authMiddleware      *middleware.AuthMiddleware
}

func NewRouter(
	clientHandler *handler.ClientHandler,
	consultationHandler *handler.ConsultationHandler,
	patientHandler *handler.PatientHandler,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		mux:                 http.NewServeMux(),
		clientHandler:       clientHandler,
		consultationHandler: consultationHandler,
		patientHandler:      patientHandler,
		userHandler:         userHandler,
		authMiddleware:      authMiddleware,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {

	//USERS
	r.mux.HandleFunc("POST /api/auth/login", r.userHandler.LogInHandler)
	r.mux.HandleFunc("POST /api/users", r.userHandler.CreateUserHandler)
	r.mux.HandleFunc("POST /api/auth/logout", r.authMiddleware.Authenticate(r.userHandler.LogOutHandler))
	r.mux.HandleFunc("DELETE /api/users/{user-id}", r.authMiddleware.Authenticate(r.userHandler.DeleteUserHandler))
	r.mux.HandleFunc("PUT /api/users/{user-id}", r.authMiddleware.Authenticate(r.userHandler.UpdateUserHandler))
	r.mux.HandleFunc("PUT /api/users/{user-id}/password", r.authMiddleware.Authenticate(r.userHandler.UpdatePasswordHandler))

	//CLIENTS
	r.mux.HandleFunc("POST /api/clients", r.authMiddleware.Authenticate(r.clientHandler.CreateClient))
	r.mux.HandleFunc("GET /api/clients/{client-id}", r.authMiddleware.Authenticate(r.clientHandler.GetClientByIDHandler))
	r.mux.HandleFunc("GET /api/clients/dni/{client-dni}", r.authMiddleware.Authenticate(r.clientHandler.GetClientByDNIHandler))
	r.mux.HandleFunc("PUT /api/clients/{client-id}", r.authMiddleware.Authenticate(r.clientHandler.UpdateClientHandler))
	r.mux.HandleFunc("DELETE /api/clients/{client-id}", r.authMiddleware.Authenticate(r.clientHandler.DeleteClientHandler))

	//PATIENTS
	r.mux.HandleFunc("POST /api/patients", r.authMiddleware.Authenticate(r.patientHandler.CreatePatientHandler))
	r.mux.HandleFunc("GET /api/patients/{patient-id}", r.authMiddleware.Authenticate(r.patientHandler.GetPatientByIDHandler))
	r.mux.HandleFunc("GET /api/patients/owner/{owner-id}", r.authMiddleware.Authenticate(r.patientHandler.GetPatientByOwnerIDHandler))
	r.mux.HandleFunc("PUT /api/patients/{patient-id}", r.authMiddleware.Authenticate(r.patientHandler.UpdatePatientHandler))
	r.mux.HandleFunc("DELETE /api/patients/{patient-id}", r.authMiddleware.Authenticate(r.patientHandler.DeletePatientHandler))

	//CONSULTATIONS
	r.mux.HandleFunc("POST /api/consultations", r.authMiddleware.Authenticate(r.consultationHandler.CreateConsultationHandler))
	r.mux.HandleFunc("GET /api/consultations/{consultation-id}", r.authMiddleware.Authenticate(r.consultationHandler.GetConsultationByIDHandler))
	r.mux.HandleFunc("GET /api/clients/{client-id}/consultations", r.authMiddleware.Authenticate(r.consultationHandler.GetConsultationsByClientIDHandler))
	r.mux.HandleFunc("GET /api/patients/{patient-id}/consultations", r.authMiddleware.Authenticate(r.consultationHandler.GetConsultationsByPatientIDHandler))
	r.mux.HandleFunc("GET /api/consultations", r.authMiddleware.Authenticate(r.consultationHandler.GetAllConsultationsHandler))
	r.mux.HandleFunc("PUT /api/consultations/{consultation-id}", r.authMiddleware.Authenticate(r.consultationHandler.UpdateConsultationHandler))
	r.mux.HandleFunc("DELETE /api/consultations/{consultation-id}", r.authMiddleware.Authenticate(r.consultationHandler.DeleteConsultationHandler))

	return r.mux
}
