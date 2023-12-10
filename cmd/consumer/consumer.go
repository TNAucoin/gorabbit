package main

import (
	"fmt"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
)

// main is the entry point of the program.
// It creates a new RabbitMQ instance and consumes messages from it.
// If any error occurs during the process, it will be handled.
// The function runs an infinite loop to continuously receive messages from the channel.
// It prints the received message and waits for more messages.
func main() {
	rabbitMQ, err := rmq.NewRabbitMQ()
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.Consume()
	errors.HandleErrorWithMessage(err, "could not consume messages")

	var forever chan struct{}
	go func() {
		for msg := range msgs {
			log.Printf("received message: %s", msg.Body)
		}
	}()
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
