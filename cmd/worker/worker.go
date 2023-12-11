package main

import (
	"bytes"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
	"time"
)

// main function is the entry point of the program.
func main() {
	rmq, err := rmq.NewRabbitMQ("task", true)
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	receivedMessages, err := rmq.Consume()
	errors.HandleErrorWithMessage(err, "could not consume messages")
	// This helps spread the work around to free worker nodes, by assigning new tasks to idle workers if possible
	err = rmq.Channel.Qos(1, 0, false)
	errors.HandleErrorWithMessage(err, "could not configure QoS")
	var forever chan struct{}
	go func() {
		for msg := range receivedMessages {
			log.Printf("received message: %s", msg.Body)
			dotCount := bytes.Count(msg.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("done")
			msg.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
