package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
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
