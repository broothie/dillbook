package main

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/config"
	"github.com/broothie/dillbook/model"
	"github.com/broothie/dillbook/server"
	_ "github.com/joho/godotenv/autoload"
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

	if err := srv.DB.AutoMigrate(
		new(model.User),
		new(model.Court),
		new(model.Booking),
	); err != nil {
		logger.Error("failed to auto migrate db", zap.Error(err))
		return
	}

	logger.Info("starting server", zap.Any("config", cfg))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), srv.Handler()); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", zap.Error(err))
	}
}
