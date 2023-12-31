package rmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQConnection creates a new RabbitMQ instance. It returns an error if any. It takes the following parameters:
// username: The username of the RabbitMQ server
// password: The password of the RabbitMQ server
// host: The host of the RabbitMQ server
// vhost: The virtual host of the RabbitMQ server
// returns: A new *amqp.connection
func NewRabbitMQConnection(username, password, host, vhost string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewRabbitMQClient creates a new RabbitMQ instance. It returns an error if any. It takes the following parameters:
// conn: The connection to the RabbitMQ server
// returns: A new RabbitMQ instance
func NewRabbitMQClient(conn *amqp.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return &RabbitClient{}, err
	}
	if err = ch.Confirm(false); err != nil {
		return &RabbitClient{}, err
	}
	return &RabbitClient{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// QueueDeclare declares a queue to hold messages and deliver to consumers.
func (rmq *RabbitClient) QueueDeclare(name string, durable, autoDelete bool) error {
	_, err := rmq.Channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// CreateBinding creates a binding between a queue and an exchange.
func (rc *RabbitClient) CreateBinding(name, binding, exchange string) error {
	return rc.Channel.QueueBind(name, binding, exchange, false, nil)
}

// Send sends a message to the queue
func (rc *RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	// Set manditory to true to return message if no queue is bound to the exchange
	confirmation, err := rc.Channel.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		routingKey,
		true,
		false,
		options,
	)
	if err != nil {
		return err
	}
	log.Println(confirmation.Wait())
	return nil
}

func (rc *RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.Channel.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rmq *RabbitClient) Close() {
	_ = rmq.Channel.Close()
}
