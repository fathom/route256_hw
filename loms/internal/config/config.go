package config

import (
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type ConfigStruct struct {
	GrpcPort     string `env:"GRPC_PORT" envDefault:"50051"`
	DatabaseURL  string `env:"DATABASE_URL" envDefault:"postgres://postgres:secret@loms-db:5432/loms?sslmode=disable"`
	Dev          bool   `env:"DEVELOPMENT_MODE" envDefault:"true"`
	Jaeger       string `env:"JAEGER"`
	KafkaBrokers []string
}

var ConfigData ConfigStruct

func Init() error {
	if err := env.Parse(&ConfigData); err != nil {
		return errors.WithMessage(err, "parsing ENV")
	}
	ConfigData.KafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	return nil
}
