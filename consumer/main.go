package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"kafka-batch-processing/lib"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)

	partitionConsumer, err := consumer.ConsumePartition(os.Getenv("KAFKA_TOPIC"), 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	isStop := false
	ringBufferSize := uint64(20)
	ringBuffer := lib.NewRingBuffer[sarama.ConsumerMessage](ringBufferSize)

	wg.Add(1)
	go func() {
		wg.Done()
		for {
			if isStop {
				break
			}

			if ringBuffer.Length() >= 10 {
				for {
					msg, ok := ringBuffer.Dequeue()
					if !ok {
						break
					}
					fmt.Printf("Received message: %s\n", msg.Value)
				}
				log.Println("stop batch processing")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			ok := ringBuffer.Enqueue(msg)
			if !ok {
				log.Println("= msg was not add into the queue", string(msg.Value))
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %s\n", err.Error())
		case <-sigterm:
			log.Println("Shutting down consumer...")
			isStop = true
			wg.Wait()
			return
		}
	}

}
