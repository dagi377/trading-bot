package api

import (
	"database/sql"
	"log"
	"net/http"
)

// Server represents the API server
type Server struct {
	port string
	db   *sql.DB
	auth *AuthService
}

// NewServer creates a new API server
func NewServer(port string, db *sql.DB) *Server {
	return &Server{
		port: port,
		db:   db,
		auth: NewAuthService(db),
	}
}

// Start starts the API server
func (s *Server) Start() error {
	// Set up routes
	http.HandleFunc("/api/login", s.auth.LoginHandler)

	// Protected routes
	http.HandleFunc("/api/protected", s.auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Protected endpoint"))
	}))

	log.Printf("Starting API server on port %s", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
