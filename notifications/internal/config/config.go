package config

import (
	"os"
	"strings"
)

type ConfigStruct struct {
	KafkaBrokers []string
}

var ConfigData ConfigStruct

func Init() {
	ConfigData.KafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
}
