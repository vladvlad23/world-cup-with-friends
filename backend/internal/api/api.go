package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"worldcup/backend/internal/auth"
)

type Server struct {
	pool *pgxpool.Pool
	auth *auth.Service
}

func NewServer(pool *pgxpool.Pool, authSvc *auth.Service) *Server {
	return &Server{pool: pool, auth: authSvc}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/login", s.handleLogin)

		r.Group(func(r chi.Router) {
			r.Use(s.auth.Middleware)
			r.Get("/me", s.handleMe)
			r.Post("/me/password", s.handleChangePassword)
			r.Get("/users", s.handleUsers)
			r.Get("/matches", s.handleMatches)
			r.Get("/matches/{id}", s.handleMatch)
			r.Post("/matches/{id}/meetups", s.handleCreateMeetup)
			r.Get("/meetups/{id}", s.handleMeetup)
			r.Post("/meetups/{id}/join", s.handleJoinMeetup)
		})
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
