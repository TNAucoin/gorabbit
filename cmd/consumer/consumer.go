package main

import (
	"fmt"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
)

func main() {
	rabbitMQ, err := rmq.NewRabbitMQ()
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	defer rabbitMQ.Close()
	q, err := rabbitMQ.Channel.QueueDeclare(
		"gorabbit",
		false, //durable
		false, //delete when unused
		false, //exclusive
		false, //no-wait
		nil,   //arguments
	)
	errors.HandleErrorWithMessage(err, "could not declare queue")
	msgs, err := rabbitMQ.Channel.Consume(q.Name, "", true, false, false, false, nil)
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
