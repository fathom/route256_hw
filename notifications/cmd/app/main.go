package main

import (
	"context"
	"log"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	lg "route256/notifications/internal/logger"
	OrderStatus "route256/notifications/internal/notifications/order_status"

	"go.uber.org/zap"
)

// Notifications
// Будет слушать Кафку и отправлять уведомления, внешнего API нет.

func main() {
	config.Init()

	logger := lg.NewLogger(config.ConfigData.Dev)
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatal("logger sync", err)
		}
	}()

	consumer, err := kafka.NewConsumer(config.ConfigData.KafkaBrokers)
	if err != nil {
		logger.Fatal("unable to create kafka consumer", zap.Error(err))
	}

	receiver := OrderStatus.NewReceiver(consumer)
	err = receiver.Subscribe("orders")
	if err != nil {
		logger.Fatal("get error from kafka", zap.Error(err))
	}

	<-context.TODO().Done()
}
