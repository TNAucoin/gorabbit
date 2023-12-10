package rmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tnaucoin/gorabbit/internal/data"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"gorabbit",
		false, //durable
		false, //delete when unused
		false, //exclusive
		false, //no-wait
		nil,   //arguments
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{Connection: conn, Channel: ch, Queue: q}, nil
}

// Publish publishes a message to the queue
func (rmq *RabbitMQ) Publish(ctx context.Context, data data.MyData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = rmq.Channel.PublishWithContext(ctx, "", rmq.Queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
	})
	return err
}

func (rmq *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return rmq.Channel.Consume(rmq.Queue.Name, "", false, false, false, false, nil)
}

func (rmq *RabbitMQ) Close() {
	_ = rmq.Channel.Close()
	_ = rmq.Connection.Close()
}
