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
	// Set up a signal handler to catch Ctrl+C and cleanly exit the program
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Set up a Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Connect to the Kafka brokers
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

	// Send a message to the "ABC" topic every 30 seconds
	xx := 0
	for {
		select {
		case <-time.After(5 * time.Second):
			for i := 0; i < 10; i++ {
				message := &sarama.ProducerMessage{
					Topic: os.Getenv("KAFKA_TOPIC"),
					Value: sarama.StringEncoder(fmt.Sprintf("Hello %d, Kafka!", xx)),
				}
				partition, offset, err := producer.SendMessage(message)
				if err != nil {
					fmt.Printf("Failed to send message: %v\n", err)
				} else {
					fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
				}
				xx += 1
			}
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			return
		}
	}
}
