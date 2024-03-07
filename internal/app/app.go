package app

import (
	"github.com/arxon31/bank/config"
	mq "github.com/arxon31/bank/internal/controller/rabbtimq"
	"github.com/arxon31/bank/internal/usecase"
	"github.com/arxon31/bank/internal/usecase/repo/transaction"
	"github.com/arxon31/bank/internal/usecase/repo/user"

	"github.com/arxon31/bank/pkg/amqp"
	"github.com/arxon31/bank/pkg/postgres"

	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	const op = "app.Run"
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(handler).With(slog.String("op", op))

	db, err := postgres.New(cfg.DB.URL, logger)
	if err != nil {
		logger.Error("failed to connect to postgres", err)
	}
	defer db.Close()

	transactionRepo, err := transaction.NewRepo(db)
	if err != nil {
		logger.Error("failed to create transaction repo", err)
		return
	}
	userRepo, err := user.NewRepo(db)
	if err != nil {
		logger.Error("failed to create user repo", err)
		return
	}

	rabbitConn, err := amqp.NewRabbitMQConn(cfg.AMQP.URL)
	if err != nil {
		logger.Error("failed to connect to rabbitmq", err)
		return
	}
	defer rabbitConn.Close()

	useCase := usecase.NewTransactionUseCase(transactionRepo, userRepo, logger)

	consumer := mq.NewTransactionConsumer(rabbitConn, logger, useCase)

	go func() {
		err = consumer.StartConsumer(
			cfg.AMQP.WorkerPoolSize,
			cfg.AMQP.Exchange,
			cfg.AMQP.Queue,
			cfg.AMQP.RoutingKey,
			cfg.AMQP.ConsumerTag,
		)

	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("received signal", slog.String("signal", s.String()))
	}

}
