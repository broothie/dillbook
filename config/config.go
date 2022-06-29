package config

import (
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Environment        string `envconfig:"ENVIRONMENT" required:"true"`
	Port               int    `envconfig:"PORT" default:"8080"`
	DBConnectionString string `envconfig:"DB_CONNECTION_STRING" required:"true" json:"-"`
}

func New() (*Config, error) {
	var config Config
	if err := envconfig.Process("dillbook", &config); err != nil {
		return nil, errors.Wrap(err, "failed to process config")
	}

	return &config, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) GormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				Colorful:      c.IsDevelopment(),
				LogLevel:      logger.Info,
			},
		),
	}
}
