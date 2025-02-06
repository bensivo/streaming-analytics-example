package main

import (
	"github.com/bensivo/streaming-analytics-example/webapp/kafka"
	"github.com/bensivo/streaming-analytics-example/webapp/service"
	"github.com/bensivo/streaming-analytics-example/webapp/controller"
)


func main() {
	// Connect to kafka
	kafkaClient := kafka.NewKafkaClient("kafka:9092")
	kafkaClient.Start()

	service := service.NewService(kafkaClient)

	controller := controller.NewHttpController(service)
	controller.Start()
}
