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
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

func NewRouter(
	clientHandler *handler.ClientHandler,
	consultationHandler *handler.ConsultationHandler,
	patientHandler *handler.PatientHandler,
	userHandler *handler.UserHandler,
) *Router {
	return &Router{
		mux:                 http.NewServeMux(),
		clientHandler:       clientHandler,
		consultationHandler: consultationHandler,
		patientHandler:      patientHandler,
		userHandler:         userHandler,
		authMiddleware:      &middleware.AuthMiddleware{SessionRepo: userHandler.SessionRepo},
		rateLimitMiddleware: middleware.NewRateLimitMiddleware(),
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	r.mux.Handle("GET /", http.FileServer(http.Dir("web")))
	r.mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	r.mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/register.html")
	})
	r.mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/login.html")
	})

	//USERS
	r.mux.HandleFunc("POST /api/auth/login", r.rateLimitMiddleware.RateLimit(r.userHandler.LogInHandler))
	r.mux.HandleFunc("POST /api/users", r.rateLimitMiddleware.RateLimit(r.userHandler.CreateUserHandler))
	r.mux.HandleFunc("POST /api/auth/logout", r.authMiddleware.Authenticate(r.userHandler.LogOutHandler))
	r.mux.HandleFunc("DELETE /api/users/{user_id}", r.authMiddleware.Authenticate(r.userHandler.DeleteUserHandler))
	r.mux.HandleFunc("PUT /api/users/{user_id}", r.authMiddleware.Authenticate(r.userHandler.UpdateUserHandler))
	r.mux.HandleFunc("PUT /api/users/{user_id}/password", r.authMiddleware.Authenticate(r.userHandler.UpdatePasswordHandler))

	//CLIENTS
	r.mux.HandleFunc("POST /api/clients", r.authMiddleware.Authenticate(r.clientHandler.CreateClient))
	r.mux.HandleFunc("GET /api/clients/{client_id}", r.authMiddleware.Authenticate(r.clientHandler.GetClientByIDHandler))
	r.mux.HandleFunc("GET /api/clients/dni/{client_dni}", r.authMiddleware.Authenticate(r.clientHandler.GetClientByDNIHandler))
	r.mux.HandleFunc("PUT /api/clients/{client_id}", r.authMiddleware.Authenticate(r.clientHandler.UpdateClientHandler))
	r.mux.HandleFunc("DELETE /api/clients/{client_id}", r.authMiddleware.Authenticate(r.clientHandler.DeleteClientHandler))

	//PATIENTS
	r.mux.HandleFunc("POST /api/patients", r.authMiddleware.Authenticate(r.patientHandler.CreatePatientHandler))
	r.mux.HandleFunc("GET /api/patients/{patient_id}", r.authMiddleware.Authenticate(r.patientHandler.GetPatientByIDHandler))
	r.mux.HandleFunc("GET /api/patients/owner/{owner_id}", r.authMiddleware.Authenticate(r.patientHandler.GetPatientByOwnerIDHandler))
	r.mux.HandleFunc("PUT /api/patients/{patient_id}", r.authMiddleware.Authenticate(r.patientHandler.UpdatePatientHandler))
	r.mux.HandleFunc("DELETE /api/patients/{patient_id}", r.authMiddleware.Authenticate(r.patientHandler.DeletePatientHandler))

	//CONSULTATIONS
	r.mux.HandleFunc("POST /api/consultations", r.authMiddleware.Authenticate(r.consultationHandler.CreateConsultationHandler))
	r.mux.HandleFunc("GET /api/consultations/{consultation_id}", r.authMiddleware.Authenticate(r.consultationHandler.GetConsultationByIDHandler))
	r.mux.HandleFunc("GET /api/clients/consultations/{client_id}", r.authMiddleware.Authenticate(r.consultationHandler.GetConsultationsByClientIDHandler))
	r.mux.HandleFunc("GET /api/patients/consultations/{patient_id}", r.authMiddleware.Authenticate(r.consultationHandler.GetConsultationsByPatientIDHandler))
	r.mux.HandleFunc("GET /api/consultations", r.authMiddleware.Authenticate(r.consultationHandler.GetAllConsultationsHandler))
	r.mux.HandleFunc("PUT /api/consultations/{consultation_id}", r.authMiddleware.Authenticate(r.consultationHandler.UpdateConsultationHandler))
	r.mux.HandleFunc("DELETE /api/consultations/{consultation_id}", r.authMiddleware.Authenticate(r.consultationHandler.DeleteConsultationHandler))

	return r.mux
}
