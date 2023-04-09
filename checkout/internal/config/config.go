package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Token        string         `yaml:"token" env:"TOKEN"`
	GrpcPort     string         `yaml:"grpc_port" env:"GRPC_PORT"`
	DatabaseURL  string         `yaml:"database_url" env:"DATABASE_URL"`
	CountWorkers int            `yaml:"count_workers" env:"COUNT_WORKERS"`
	Services     ConfigServices `yaml:"services"`
	Dev          bool           `yaml:"development_mode" env:"DEVELOPMENT_MODE"`
}

type ConfigServices struct {
	Loms    string `yaml:"loms" env:"LOMS_URL"`
	Product string `yaml:"product" env:"PRODUCT_URL"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	if err := env.Parse(&ConfigData); err != nil {
		return errors.WithMessage(err, "parsing ENV")
	}

	return nil
}
