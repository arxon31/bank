package rabbtimq

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/arxon31/bank/internal/entity"
	"github.com/arxon31/bank/internal/usecase"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log/slog"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type TransactionConsumer struct {
	amqpConn           *amqp.Connection
	logger             *slog.Logger
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionConsumer(amqpConn *amqp.Connection, logger *slog.Logger, transactionUseCase usecase.TransactionUseCase) *TransactionConsumer {
	return &TransactionConsumer{
		amqpConn:           amqpConn,
		logger:             logger,
		transactionUseCase: transactionUseCase,
	}
}

func (c *TransactionConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	c.logger.Info("Declaring exchange", slog.String("exchange", exchangeName))
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	c.logger.Info("Declared queue, binding it to exchange",
		slog.String("queue", queue.Name),
		slog.Int("messages count", queue.Messages),
		slog.Int("consumer count", queue.Consumers),
		slog.String("exchange", exchangeName),
		slog.String("bindingKey", bindingKey),
	)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	c.logger.Info("Queue bound to exchange, starting to consume from queue",
		slog.String("consumer tag", consumerTag))

	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

func (c *TransactionConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {

	for delivery := range messages {
		model := entity.Transaction{}

		err := json.Unmarshal(delivery.Body, &model)
		if err != nil {
			c.logger.Error("Failed to unmarshal delivery", slog.String("error", err.Error()))
			delivery.Reject(false)
			continue
		}

		err = model.Validate()
		if err != nil {
			c.logger.Error("Failed to validate model", slog.String("error", err.Error()))
			delivery.Reject(false)
			continue
		}

		_, err = c.transactionUseCase.MakeTransaction(ctx, model)
		if err != nil {
			switch err {
			case usecase.ErrInsufficientFunds:
				c.logger.Error("Failed to make transaction", slog.String("error", err.Error()))
				delivery.Reject(false)
			case sql.ErrNoRows:
				c.logger.Error("Failed to make transaction", slog.String("error", err.Error()))
				delivery.Reject(false)
			default:
				c.logger.Error("Failed to make transaction", slog.String("error", err.Error()))
				delivery.Reject(true)
			}
		}

		c.logger.Info("Transaction processed")
		delivery.Ack(false)
	}

	c.logger.Info("Deliveries channel closed")
}

// StartConsumer Start new rabbitmq consumer
func (c *TransactionConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	//c.logger.Error("ch.NotifyClose", slog.String("error", chanErr.Error()))
	return chanErr
}
