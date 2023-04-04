package main

import (
	"context"
	"log"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	OrderStatus "route256/notifications/internal/notifications/order_status"
)

// Notifications
// Будет слушать Кафку и отправлять уведомления, внешнего API нет.

func main() {
	config.Init()

	consumer, err := kafka.NewConsumer(config.ConfigData.KafkaBrokers)
	if err != nil {
		log.Fatalf("Unable to create kafka consumer: %v", err)
	}

	receiver := OrderStatus.NewReceiver(consumer)
	err = receiver.Subscribe("orders")
	if err != nil {
		log.Fatalf("Get error from kafka: %v", err)
	}

	<-context.TODO().Done()
}
