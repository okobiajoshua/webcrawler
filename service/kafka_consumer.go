package service

import (
	"context"
	"fmt"

	"github.com/monzo/webcrawler/data"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// KafkaConsumer struct
type KafkaConsumer struct {
	cons []*kafka.Consumer
}

// NewConsumer returns a Consumer struct
func NewConsumer(data data.Data) *KafkaConsumer {
	kc := &KafkaConsumer{cons: []*kafka.Consumer{}}
	for i := int32(0); i < 10; i++ {
		con, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "kafka:9092",
			"group.id":          "myGroup",
			"auto.offset.reset": "earliest",
		})

		if err != nil {
			panic(err)
		}

		con.Assign([]kafka.TopicPartition{{Topic: &topic, Partition: i}})
		kc.cons = append(kc.cons, con)
	}

	return kc
}

// Consume method processes streams
func (c *KafkaConsumer) Consume(ctx context.Context, f func([]byte) error) {
	for _, con := range c.cons {
		go listen(con, f)
	}
}

func listen(con *kafka.Consumer, f func([]byte) error) {
	defer con.Close()
	for {
		msg, err := con.ReadMessage(-1)
		if err == nil {
			f(msg.Value)
			// log.Println("FROM PARTITION:- ", string(msg.Value), msg.TopicPartition)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
