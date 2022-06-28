package application

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

type Application struct {
	DB     *gorm.DB
	logger *zap.Logger
	config *config.Config
	render *render.Render
	conn   *sql.DB
}

func New(cfg *config.Config, logger *zap.Logger) (*Application, error) {
	conn, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sql connection")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open gorm")
	}

	return &Application{
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

func (a *Application) Close() error {
	return a.conn.Close()
}
