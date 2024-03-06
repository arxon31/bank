package amqp

import "github.com/streadway/amqp"

type AMQP struct {
	Channel *amqp.Channel
}

func New(url string, queueName string) (*AMQP, error) {
	connectRabbitMQ, err := amqp.Dial(url)
	if err != nil {
		return &AMQP{}, err
	}

	channel, err := connectRabbitMQ.Channel()
	if err != nil {
		return &AMQP{}, err
	}

	_, err = channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return &AMQP{}, err
	}

	return &AMQP{
		Channel: channel,
	}, nil

}

func (a *AMQP) Close() error {
	return a.Channel.Close()
}
