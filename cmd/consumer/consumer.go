package main

import (
	"context"
	"fmt"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

// main is the entry point of the program.
// It creates a new RabbitMQ instance and consumes messages from it.
// If any error occurs during the process, it will be handled.
// The function runs an infinite loop to continuously receive messages from the channel.
// It prints the received message and waits for more messages.
func main() {
	conn, err := rmq.NewRabbitMQConnection("guest", "guest", "localhost:5672", "gorabbit")
	errors.HandleErrorWithMessage(err, "could not create rabbitmq connection")
	client, err := rmq.NewRabbitMQClient(conn)
	errors.HandleErrorWithMessage(err, "could not create rabbitmq client")
	msgBus, err := client.Consume("gorabbit", "gorabbit_consumer", false)
	errors.HandleErrorWithMessage(err, "could not consume messages")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	//Errgroup to handle multiple goroutines
	g, ctx := errgroup.WithContext(ctx)
	// Limit the number of goroutines to 10
	g.SetLimit(10)
	var forever chan struct{}
	go func() {
		for message := range msgBus {
			msg := message
			g.Go(func() error {
				log.Printf("received message: %v", msg)
				time.Sleep(time.Second * 10)
				if err := msg.Ack(false); err != nil {
					fmt.Println("acknoledge message failed: retry ? handle manually %s\n", msg.MessageId)
					return err
				}
				log.Printf("message acknowledged %s\n", msg.MessageId)
				return nil
			})
		}
	}()
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
