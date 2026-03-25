package main

import (
	"context"
	"log"
	"time"

	"github.com/iamonah/eventdriven/internal"
	"github.com/rabbitmq/amqp091-go"
)

// we have the customers_events exchange created on via the cli
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

	if err := client.CreateQueue("customer_created", true, false); err != nil {
		panic(err)
	}
	if err := client.CreateQueue("customer_test", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customer_created", "customers.created.*", "customer_events"); err != nil {
		panic(err)
	}
	if err := client.CreateBinding("customer_test", "customers.*", "customer_events"); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		if err := client.Send(ctx, "customer_events", "customers.created.us", amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         []byte("A cool message between services"),
		}); err != nil {
			panic(err)
		}
	}
	//sending a transient message
	if err := client.Send(ctx, "customer_events", "customers.test", amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Transient,
		Body:         []byte("uncool undurable message"),
	}); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 60)

	log.Println(client)
}
