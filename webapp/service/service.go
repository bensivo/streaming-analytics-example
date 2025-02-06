package service

import (
	"github.com/bensivo/streaming-analytics-example/webapp/kafka"
)

type Service struct {
	kafkaClient *kafka.KafkaClient
}

func NewService(kafkaClient *kafka.KafkaClient) *Service {
	return &Service{
		kafkaClient: kafkaClient,
	}
}

type Rating struct {
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
	Rating    int    `json:"rating"`
}

func (s *Service) SendRating(rating Rating) {
	s.kafkaClient.ProduceJson("event.customer_feedback", rating)
}