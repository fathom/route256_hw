package order_status

import (
	"fmt"
	"route256/notifications/internal/logger"

	"github.com/Shopify/sarama"
)

type OrderStatusReceiver interface {
	Subscribe(topic string) error
}

type receiver struct {
	consumer sarama.Consumer
}

func NewReceiver(consumer sarama.Consumer) *receiver {
	return &receiver{
		consumer: consumer,
	}
}

func (r *receiver) Subscribe(topic string) error {
	partitionList, err := r.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		pc, err := r.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				orderID := string(message.Key)
				status := string(message.Value)

				logger.Info(fmt.Sprintf(
					"read: orderID: %s, status: %s,  topic: %s, partion: %d, offset: %d",
					orderID,
					status,
					topic,
					message.Partition,
					message.Offset,
				))
			}
		}(pc)
	}

	return nil
}
