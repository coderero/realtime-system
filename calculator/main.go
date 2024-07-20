package main

import (
	"github.com/sirupsen/logrus"
)

const (
	kafkaTopic = "obu_data"
)

func main() {
	calculatorService := NewCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, calculatorService)
	if err != nil {
		logrus.Fatalf("Error creating Kafka consumer: %v", err)
	}

	kafkaConsumer.Start()
}
