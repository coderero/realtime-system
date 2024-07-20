package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/coderero/toll_calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	Produce(types.OBUData) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    *string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &KafkaProducer{producer: p, topic: &topic}, nil
}

func (kp *KafkaProducer) Produce(data types.OBUData) error {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(data); err != nil {
		return err
	}

	return kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: kp.topic, Partition: kafka.PartitionAny},
		Value:          buff.Bytes(),
	}, nil)
}
