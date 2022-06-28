package server

import (
	"database/sql"

	"github.com/broothie/dillbook/config"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	logger *zap.Logger
	config *config.Config
	render *render.Render
	conn   *sql.DB
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
		logger: logger,
		config: cfg,
		render: render.New(render.Options{
			Layout:        "layout",
			Extensions:    []string{".gohtml"},
			IsDevelopment: cfg.IsDevelopment(),
		}),
		conn: conn,
		DB:   db,
	}, nil
}

func (s *Server) Close() error {
	return s.conn.Close()
}
