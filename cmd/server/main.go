package main

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/config"
	"github.com/broothie/dillbook/model"
	"github.com/broothie/dillbook/server"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger
	if cfg.IsDevelopment() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println("failed to sync logger", err)
		}
	}()

	srv, err := server.New(cfg, logger)
	if err != nil {
		logger.Error("failed to create application", zap.Error(err))
		return
	}
	defer func() {
		if err := srv.Close(); err != nil {
			logger.Error("failed to close application", zap.Error(err))
		}
	}()

	if err := initDB(srv, logger); err != nil {
		logger.Error("failed to initialize db", zap.Error(err))
		return
	}

	logger.Info("starting server", zap.Any("config", cfg))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), srv.Handler()); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", zap.Error(err))
	}
}

func initDB(srv *server.Server, logger *zap.Logger) error {
	logger = logger.With(zap.String("at", "initDB"))

	logger.Info("auto migrating")
	if err := srv.DB.AutoMigrate(
		new(model.Location),
		new(model.Court),
		new(model.Booking),
	); err != nil {
		return errors.Wrap(err, "failed to auto migrate db")
	}

	logger.Info("creating extensions")
	if err := srv.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return errors.Wrap(err, "failed to create uuid-ossp db extension")
	}

	return nil
}
