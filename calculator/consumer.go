package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/coderero/toll_calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type kafkaConsumer struct {
	consumer   *kafka.Consumer
	isRunning  bool
	calculator CalculatorService
}

func NewKafkaConsumer(topic string, s CalculatorService) (*kafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		consumer:   c,
		calculator: s,
	}, nil
}

func (k *kafkaConsumer) Start() {
	logrus.Info("Starting Kafka consumer")
	k.isRunning = true
	k.Consume()
}

func (k *kafkaConsumer) Consume() {
	for k.isRunning {
		var obuData types.OBUData
		msg, err := k.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Error reading message: %v", err)
			continue
		}

		reader := bytes.NewReader(msg.Value)
		if err := json.NewDecoder(reader).Decode(&obuData); err != nil {
			logrus.Errorf("Error decoding message: %v", err)
			continue
		}

		fmt.Printf("Received message: %v\n", obuData)

	}
}
