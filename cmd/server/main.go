package main

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/application"
	"github.com/broothie/dillbook/config"
	"github.com/broothie/dillbook/model"
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

	app, err := application.New(cfg, logger)
	if err != nil {
		logger.Error("failed to create application", zap.Error(err))
		return
	}
	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("failed to close application", zap.Error(err))
		}
	}()

	if err := app.DB.AutoMigrate(
		new(model.User),
		new(model.Court),
		new(model.Booking),
	); err != nil {
		logger.Error("failed to auto migrate db", zap.Error(err))
		return
	}

	logger.Info("starting server", zap.Any("config", cfg))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), app.Handler()); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", zap.Error(err))
	}
}
