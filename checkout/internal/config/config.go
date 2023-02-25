package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Token    string         `yaml:"token" env:"TOKEN"`
	AppPort  string         `yaml:"app_port" env:"APP_PORT"`
	Services ConfigServices `yaml:"services"`
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

	log.Printf("%+v", ConfigData.Services)

	return nil
}