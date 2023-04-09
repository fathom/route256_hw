package order_status

import (
	"fmt"
	"route256/loms/internal/logger"
	"route256/loms/internal/model"
	"time"

	"github.com/Shopify/sarama"
)

type OrderStatusSender interface {
	SendOrderStatus(orderID int64, status model.OrderStatus) error
}

type orderStatusSender struct {
	producer sarama.SyncProducer
	topic    string
}

type Handler func(id string)

func NewOrderStatusSender(producer sarama.SyncProducer, topic string) OrderStatusSender {
	return orderStatusSender{
		producer: producer,
		topic:    topic,
	}
}

func (o orderStatusSender) SendOrderStatus(orderID int64, status model.OrderStatus) error {
	msg := &sarama.ProducerMessage{
		Topic:     o.topic,
		Key:       sarama.StringEncoder(fmt.Sprint(orderID)),
		Value:     sarama.StringEncoder(status),
		Partition: -1,
		Timestamp: time.Now(),
	}

	partition, offset, err := o.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("order id: %d, status: %v, partition: %d, offset: %d", orderID, status, partition, offset))

	return nil
}
