package main

import (
	"log"
	"time"

	"github.com/iamonah/eventdriven/internal"
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
	time.Sleep(time.Second * 60)
}
