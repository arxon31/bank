package rabbtimq

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"time"

	"github.com/arxon31/bank/internal/entity"
	"github.com/streadway/amqp"
)

const tickerDuration = 3 * time.Second

type TransactionPublisher struct {
	logger  *slog.Logger
	channel *amqp.Channel

	exchangeName string
	routingKey   string
}

func NewTransactionPublisher(amqpConn *amqp.Connection, logger *slog.Logger, exchangeName, routingKey string) *TransactionPublisher {
	ch, err := amqpConn.Channel()
	if err != nil {
		logger.Error("amqp connection error", slog.String("error", err.Error()))
		return &TransactionPublisher{}
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		logger.Error("Failed to declare an exchange", slog.String("error", err.Error()))
		return &TransactionPublisher{}
	}

	return &TransactionPublisher{
		logger:       logger,
		channel:      ch,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}
}

func (p *TransactionPublisher) Start(ctx context.Context) {
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			transaction, err := p.generateTransactionJSON()
			if err != nil {
				p.logger.Error("Failed to generate transaction", slog.String("error", err.Error()))
				continue
			}
			p.publishMessage(transaction)
			p.logger.Info("Published message", slog.String("message", string(transaction)))
		case <-ctx.Done():
			p.stop()
		}
	}

}

func (p *TransactionPublisher) publishMessage(body []byte) {

	err := p.channel.Publish(
		p.exchangeName,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		p.logger.Error("Failed to publish message", slog.String("error", err.Error()))
	}
}

func (p *TransactionPublisher) generateTransactionJSON() ([]byte, error) {

	from, to := rand.Int63n(10), rand.Int63n(10)
	amount := rand.Int63n(500)

	transaction := entity.Transaction{
		FromAccountID: from,
		ToAccountID:   to,
		Amount:        amount,
	}

	return json.Marshal(transaction)
}

func (p *TransactionPublisher) stop() {
	err := p.channel.Close()
	if err != nil {
		p.logger.Error("Failed to close channel", slog.String("error", err.Error()))
	}
}
