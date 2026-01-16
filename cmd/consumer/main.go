package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iamonah/eventdriven/internal"
	"golang.org/x/sync/errgroup"
)

func main() {
	conn, err := internal.ConnectRabbitMq("user", "password", "localhost", "customers", 5672)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMqClient(conn)
	if err != nil {
		log.Fatal("failed to open client connection")
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
	go func() {
		for delivery := range payload {
			msg := delivery
			g.Go(func() error {
				fmt.Println(string(msg.Body))

				//if a  process fails
				//return error
				if !msg.Redelivered {
					msg.Nack(false, true) // or false
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
	}()

	go func() {
		for data := range message {
			msg := data
			g.Go(func() error {
				fmt.Println(string(msg.Body))

				//imagine a  process fails
				//return error
				if !msg.Redelivered {
					msg.Nack(false, true) // or false
					return nil
				}

				if err := msg.Ack(false); err != nil {
					fmt.Printf("ack failed, channel likely closed: %v", err)
					return err
				}
				log.Printf("Acknowledge message %s\n", msg.MessageId)
				return nil
			})
		}
	}()

	log.Println("consuming, to close the program press CTRL+C")

	<-blocking
}
