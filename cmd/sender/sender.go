package main

import (
	"context"
	"github.com/tnaucoin/gorabbit/internal/data"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
	"sync"
	"time"
)

// main is the entry point of the application.
// It creates a new RabbitMQ instance, initializes a WaitGroup for concurrency management,
// and starts a goroutine to publish messages to RabbitMQ.
// After the goroutine finishes, it waits for all goroutines to complete using the WaitGroup.
// Finally, it closes the RabbitMQ connection.
func main() {
	// ...
	rabbitMQ, err := rmq.NewRabbitMQ("gorabbit", false)
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for i := 0; i < 10; i++ {
			data := data.MyData{
				Message: "Hello World",
				Number:  i + 1,
			}
			err := rabbitMQ.PublishJSON(ctx, data)
			if err != nil {
				log.Fatalf("could not publish message: %v", err)
			} else {
				log.Printf("published message: %v", data)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()
	defer rabbitMQ.Close()
}
