package config

import (
	"os"
	"strings"
)

type ConfigStruct struct {
	KafkaBrokers []string
	Dev          bool
}

var ConfigData ConfigStruct

func Init() {
	ConfigData.KafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	ConfigData.Dev = os.Getenv("DEVELOPMENT_MODE") == "true"
}
