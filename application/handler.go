package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

func (a *Application) Handler() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(a.zapLoggerMiddleware)

	// Routes
	r.Get("/", a.Index)

	r.Route("/courts", func(r chi.Router) {
		r.Get("/new", a.NewCourt)
		r.Post("/", a.CreateCourt)

		r.Route("/{courtID}", func(r chi.Router) {
			r.Get("/", a.ShowCourt)
		})
	})

	return r
}

func (a *Application) zapLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(zapctx.WithLogger(r.Context(), a.logger.With(
			zap.String("request_id", r.Header.Get(middleware.RequestIDHeader)),
		))))
	})
}
