package sender

import (
	"github.com/tnaucoin/gorabbit/internal/data"
	"github.com/tnaucoin/gorabbit/internal/errors"
	"github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
	"time"
)

func main() {
	// ...
	rabbitMQ, err := rmq.NewRabbitMQ()
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	defer rabbitMQ.Close()
	go func() {
		for i := 0; i < 10; i++ {
			data := data.MyData{
				Message: "Hello World",
				Number:  i + 1,
			}
			err := rabbitMQ.Publish(data)
			if err != nil {
				log.Fatalf("could not publish message: %v", err)
			} else {
				log.Printf("published message: %v", data)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("could not consume messages: %v", err)
	}
	for msg := range msgs {
		log.Printf("received message: %s", msg.Body)
		msg.Ack(false)
	}
}
