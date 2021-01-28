package service

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	topic = "crawler"
)

// KafkaProducer struct
type KafkaProducer struct {
	p          *kafka.Producer
	mu         sync.Mutex
	partitions chan int32
	mp         map[string]int32
}

// NewKafkaProducer returns a producer struct
func NewKafkaProducer() *KafkaProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092", "enable.idempotence": true})
	if err != nil {
		panic(err)
	}
	k := &KafkaProducer{p: p}
	k.mp = map[string]int32{}
	k.partitions = make(chan int32, 10)
	for i := int32(0); i < 10; i++ {
		k.partitions <- i
	}
	// go k.deliveryReport()
	return k
}

// Publish message to event stream
func (k *KafkaProducer) Publish(ctx context.Context, msg []byte) error {
	uri, err := url.Parse(string(msg))
	if err != nil {
		return err
	}
	if uri.Hostname() == "" {
		return fmt.Errorf("url hostname is not known")
	}
	partition, err := k.getPartition(msg)
	if err != nil {
		return err
	}
	return k.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          msg,
	}, nil)
}

// Delivery report handler for produced messages
func (k *KafkaProducer) deliveryReport() {
	for e := range k.p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

func (k *KafkaProducer) getPartition(msg []byte) (int32, error) {
	uri, err := url.Parse(string(msg))
	if err != nil {
		return int32(-1), err
	}
	if _, ok := k.mp[uri.Hostname()]; !ok {
		k.mu.Lock()
		if len(k.mp) >= 10 {
			k.mu.Unlock()
			return -1, fmt.Errorf("no available space for processing now")
		}
		n := <-k.partitions
		k.mp[uri.Hostname()] = n
		k.mu.Unlock()
		return n, nil
	}
	return k.mp[uri.Hostname()], nil
}
