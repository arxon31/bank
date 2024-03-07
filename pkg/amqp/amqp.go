package amqp

import "github.com/streadway/amqp"

func NewRabbitMQConn(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
