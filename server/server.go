package server

import (
	"database/sql"
	"net/http"

	"github.com/broothie/dillbook/config"
	"github.com/gorilla/schema"
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
}

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	conn, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sql connection")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gorm")
	}

	return &Server{
		DB:     db,
		logger: logger,
		config: cfg,
		render: render.New(render.Options{
			Layout:        "layout",
			Extensions:    []string{".gohtml"},
			IsDevelopment: cfg.IsDevelopment(),
		}),
		conn:        conn,
		formDecoder: schema.NewDecoder(),
	}, nil
}

func (s *Server) Close() error {
	return s.conn.Close()
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
