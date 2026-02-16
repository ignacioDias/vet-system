package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vetsys/internal/router"
)

type Server struct {
	httpServer *http.Server
	router     *router.Router
}

func NewServer(port string, router *router.Router) *Server {
	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Initialize() error {
	s.httpServer.Handler = s.router.SetupRoutes()

	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) StartServer(router router.Router) {

	// Canal para manejar señales de interrupción
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		if err := s.Initialize(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server started successfully")

	// Esperar señal de interrupción
	<-done
	log.Println("Server stopping...")

	// Shutdown gracefully con timeout de 30 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
