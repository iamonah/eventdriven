package internal

import (
	"fmt"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	//the connection used by the client
	conn *amqp.Connection
	// channel is used to process / send messages
	ch *amqp.Channel
}

func ConnectRabbitMq(username, passwrd, host, vhost string, port int16) (*amqp.Connection, error) {
	cleanPassword := url.QueryEscape(passwrd)
	address := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, cleanPassword, host, port, vhost)

	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq %s", err)
	}
	return conn, nil
}

func NewRabbitMqClient(conn *amqp.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close() //close the client connection
}
