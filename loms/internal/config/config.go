package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type ConfigStruct struct {
	GrpcPort string `env:"GRPC_PORT"  envDefault:"50051"`
}

var ConfigData ConfigStruct

func Init() error {
	if err := env.Parse(&ConfigData); err != nil {
		return errors.WithMessage(err, "parsing ENV")
	}

	return nil
}