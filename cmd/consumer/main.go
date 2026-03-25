package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iamonah/eventdriven/internal"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

func main() {
	client, err := internal.NewRabbitMqClient(internal.RabbitMqConfig{
		Username: "user",
		Password: "password",
		Host:     "localhost",
		Vhost:    "customers",
		Port:     5672,
	})
	if err != nil {
		log.Fatalf("failed to open client connection: %v", err)
	}
	defer client.Close()

	payload, err := client.Consume("customer_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	message, err := client.Consume("customer_test", "payment-service", false)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	var blocking = make(chan int)

	go consumeMessages(payload, g)
	go consumeMessages(message, g)

	log.Println("consuming, to close the program press CTRL+C")

	<-blocking
}

func consumeMessages(messages <-chan amqp091.Delivery, g *errgroup.Group) {
	for delivery := range messages {
		msg := delivery
		g.Go(func() error {
			//example of a work to be processed
			if _, err := fmt.Println(string(msg.Body)); err != nil {
				if msg.Redelivered {
					msg.Nack(false, false) // or false
					return nil
				}
				msg.Nack(false, true)
				return nil
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("ack failed, channel likely closed: %v", err)
				return err
			}
			fmt.Printf("Acknowledge message %s\n", msg.MessageId)
			return nil
		})
	}
}
