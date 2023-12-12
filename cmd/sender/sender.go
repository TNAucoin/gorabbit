package main

import (
	"fmt"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"time"
)

// main is the entry point of the application.
// It creates a new RabbitMQ instance, initializes a WaitGroup for concurrency management,
// and starts a goroutine to publish messages to RabbitMQ.
// After the goroutine finishes, it waits for all goroutines to complete using the WaitGroup.
// Finally, it closes the RabbitMQ connection.
func main() {
	// ...
	conn, err := rmq.NewRabbitMQConnection("guest", "guest", "localhost:5672", "gorabbit")
	errors.HandleErrorWithMessage(err, "could not create rabbitmq connection")
	defer conn.Close()
	client, err := rmq.NewRabbitMQClient(conn)
	errors.HandleErrorWithMessage(err, "could not create rabbitmq client")
	defer client.Close()

	//Create Queue
	if err = client.QueueDeclare("gorabbit", true, false); err != nil {
		errors.HandleErrorWithMessage(err, "could not create queue")
	}

	if err = client.CreateBinding("gorabbit", "gorabbit.created.*", "gorabbit"); err != nil {
		errors.HandleErrorWithMessage(err, "could not create binding")
	}

	time.Sleep(time.Second * 10)
	fmt.Println(client)
}
