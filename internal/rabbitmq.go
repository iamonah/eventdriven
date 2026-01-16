package internal

import (
	"context"
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

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	//leaving nowait false makes the channel to return error if it fails to bind
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

// publish payload onto an exchange with the given routingkey
func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	//leaving nowait false makes the channel to return error if it fails to bind

	//mandatory is used to determine if a failure should drop the message that you are sending
	//or it should return an error

	//(so if mandatory is left as false it will fire and forget if set  to true if it fails to publish the message to the exchange or the
	//queue it will return an error)
	return rc.ch.PublishWithContext(ctx, exchange, routingKey,
		// Mandatory used to determine if an error should return upon failure
		true,
		//immediate remvoed in MQ 3 deprecated feature
		false,
		options)
}

// consume is used to consume a queue
func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	//always set excluse to false and nolocal is not suppoted in raggitmq but in amqp
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}
