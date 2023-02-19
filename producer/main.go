package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Println("Failed to close Kafka producer cleanly:", err)
		}
	}()

	idx := 0
	for {
		select {
		case <-time.After(5 * time.Second):
			for i := 0; i < 10; i++ {
				message := &sarama.ProducerMessage{
					Topic: os.Getenv("KAFKA_TOPIC"),
					Value: sarama.StringEncoder(fmt.Sprintf("Hello %d, Kafka!", idx)),
				}
				partition, offset, err := producer.SendMessage(message)
				if err != nil {
					fmt.Printf("Failed to send message: %v\n", err)
				} else {
					fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
				}
				idx += 1
			}
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			return
		}
	}
}
