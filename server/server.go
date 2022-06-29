package server

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/broothie/dillbook/config"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/juju/zaputil/zapctx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB          *gorm.DB
	logger      *zap.Logger
	config      *config.Config
	render      *render.Render
	conn        *sql.DB
	formDecoder *schema.Decoder
	validate    *validator.Validate
}

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	conn, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sql connection")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), cfg.GormConfig())
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gorm")
	}

	return &Server{
		DB:     db,
		logger: logger,
		config: cfg,
		render: render.New(render.Options{
			Layout:     "layout",
			Extensions: []string{".gohtml"},
			Funcs: []template.FuncMap{{
				"divide": func(a, b int) int { return a / b },
			}},
			IsDevelopment: cfg.IsDevelopment(),
		}),
		conn:        conn,
		formDecoder: schema.NewDecoder(),
		validate:    validator.New(),
	}, nil
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/locations", http.StatusPermanentRedirect)
}

func (s *Server) zapLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(zapctx.WithLogger(r.Context(), s.logger.With(
			zap.String("request_id", middleware.GetReqID(r.Context())),
		))))
	})
}

func (s *Server) decodeForm(r *http.Request, target any) error {
	if err := r.ParseForm(); err != nil {
		return errors.Wrap(err, "failed to parse form")
	}

	if err := s.formDecoder.Decode(target, r.PostForm); err != nil {
		return errors.Wrap(err, "failed to decode form")
	}

	return nil
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	zapctx.Logger(r.Context()).Error("error", zap.Error(err), zap.Int("status", code))
	http.Error(w, err.Error(), code)
}
