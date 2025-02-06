package kafka

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/IBM/sarama"
)

type KafkaClient struct {
	broker   string
	config   *sarama.Config
	producer sarama.AsyncProducer
	consumer sarama.Consumer
}

func NewKafkaClient(broker string) *KafkaClient {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	return &KafkaClient{
		broker:   broker,
		config:   config,
		producer: nil,
		consumer: nil,
	}
}

func (k *KafkaClient) Start() error {
	producer, err := sarama.NewAsyncProducer([]string{k.broker}, k.config)
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for {
			select {
			case success := <-producer.Successes():
				log.Printf("Successfully produced message to topic %s partition %d at offset %d", 
					success.Topic, success.Partition, success.Offset)

			case err := <-producer.Errors():
				log.Printf("Failed to produce message: %v", err)

			case <-signals:
				log.Println("Received interrupt signal, shutting down producer...")
				producer.Close()
				return
			}
		}
	}()

	consumer, err := sarama.NewConsumer([]string{k.broker}, k.config)
	if err != nil {
		log.Fatal(err)
	}

	k.producer = producer
	k.consumer = consumer

	return nil
}

func (k *KafkaClient) ProduceJson(topic string, message interface{}) {
	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(""),
		Value:     sarama.ByteEncoder(messageJson),
		Timestamp: time.Now(),
	}

	log.Printf("Producing message: %s\n", messageJson)
	k.producer.Input() <- msg
}

func (k *KafkaClient) Close() {
	k.producer.Close()
	k.consumer.Close()
}
