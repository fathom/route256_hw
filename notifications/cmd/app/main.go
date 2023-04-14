package main

import (
	"context"
	"net/http"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	"route256/notifications/internal/logger"
	"route256/notifications/internal/metrics"
	OrderStatus "route256/notifications/internal/notifications/order_status"

	"go.uber.org/zap"
)

// Notifications
// Будет слушать Кафку и отправлять уведомления, внешнего API нет.

func main() {
	config.Init()

	// Инициализация логирования
	logger.Init(config.ConfigData.Dev)

	consumer, err := kafka.NewConsumer(config.ConfigData.KafkaBrokers)
	if err != nil {
		logger.Fatal("unable to create kafka consumer", zap.Error(err))
	}

	receiver := OrderStatus.NewReceiver(consumer)
	err = receiver.Subscribe("orders")
	if err != nil {
		logger.Fatal("get error from kafka", zap.Error(err))
	}

	go func() {
		http.Handle("/metrics", metrics.NewMetricHandler())
		err := http.ListenAndServe(":8083", nil)
		logger.Fatal("failed metrics", zap.Error(err))
		panic(err)
	}()

	<-context.TODO().Done()
}
