package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r.Get("/", s.root)

	r.Route("/locations", func(r chi.Router) {
		r.Get("/", s.indexLocations)
		r.Get("/new", s.newLocation)
		r.Post("/", s.createLocation)

		r.Route("/{locationID}", func(r chi.Router) {
			r.Get("/", s.showLocation)

			r.Route("/courts", func(r chi.Router) {
				r.Get("/new", s.newCourt)
				r.Post("/", s.createCourt)
			})
		})
	})

	r.Route("/courts", func(r chi.Router) {
		r.Route("/{courtID}", func(r chi.Router) {
			r.Get("/", s.showCourt)
		})
	})

	return r
}
