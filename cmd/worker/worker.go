package main

import (
	"bytes"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
	"time"
)

func main() {
	rmq, err := rmq.NewRabbitMQ("task", true)
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	receivedMessages, err := rmq.Consume()
	errors.HandleErrorWithMessage(err, "could not consume messages")

	var forever chan struct{}
	go func() {
		for msg := range receivedMessages {
			log.Printf("received message: %s", msg.Body)
			dotCount := bytes.Count(msg.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("done")
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
