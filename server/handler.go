package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.zapLoggerMiddleware)

	// Routes
	r.Get("/", s.Index)

	r.Route("/users", func(r chi.Router) {
		r.Get("/new", s.NewUser)
		r.Post("/", s.CreateUser)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", s.ShowUser)
		})
	})

	r.Route("/courts", func(r chi.Router) {
		r.Get("/new", s.NewCourt)
		r.Post("/", s.CreateCourt)

		r.Route("/{courtID}", func(r chi.Router) {
			r.Get("/", s.ShowCourt)
		})
	})

	return r
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/courts/new", http.StatusPermanentRedirect)
}

func (s *Server) zapLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(zapctx.WithLogger(r.Context(), s.logger.With(
			zap.String("request_id", r.Header.Get(middleware.RequestIDHeader)),
		))))
	})
}
