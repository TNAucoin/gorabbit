package main

import (
	"context"
	"github.com/tnaucoin/gorabbit/internal/errors"
	rmq2 "github.com/tnaucoin/gorabbit/pkg/rmq"
	"log"
	"os"
	"strings"
	"time"
)

// bodyFrom returns a string from the command line arguments.
func bodyFrom(arg []string) string {
	var s string
	if (len(arg) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(arg[1:], " ")
	}
	return s
}

func main() {
	rmq, err := rmq2.NewRabbitMQ("task", true)
	errors.HandleErrorWithMessage(err, "could not create rabbitmq")
	body := bodyFrom(os.Args)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	errors.HandleErrorWithMessage(err, "could not create context")
	err = rmq.PublishString(ctx, []byte(body))
	errors.HandleErrorWithMessage(err, "could not publish message")
	log.Printf(" [x] Sent %s", body)
}
