package app

import (
	"github.com/arxon31/bank/config/cons"
	mq "github.com/arxon31/bank/internal/controller/rabbtimq"
	"github.com/arxon31/bank/internal/usecase"
	"github.com/arxon31/bank/internal/usecase/repo/transaction"
	"github.com/arxon31/bank/internal/usecase/repo/user"
	"github.com/arxon31/bank/pkg/logging"

	"github.com/arxon31/bank/pkg/amqp"
	"github.com/arxon31/bank/pkg/postgres"

	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func RunConsumer(cfg *cons.Config) {

	logger := logging.New(cfg.App.Env)

	db, err := postgres.New(cfg.DB.URL, logger)
	if err != nil {
		logger.Error("failed to connect to postgres", slog.String("url", cfg.DB.URL), slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	transactionRepo, err := transaction.NewRepo(db)
	if err != nil {
		logger.Error("failed to create transaction repo", slog.String("error", err.Error()))
		return
	}
	userRepo, err := user.NewRepo(db)
	if err != nil {
		logger.Error("failed to create user repo", slog.String("error", err.Error()))
		return
	}

	rabbitConn, err := amqp.NewRabbitMQConn(cfg.AMQP.URL)
	if err != nil {
		logger.Error("failed to connect to rabbitmq", slog.String("url", cfg.AMQP.URL), slog.String("error", err.Error()))
		return
	}
	defer rabbitConn.Close()

	useCase := usecase.NewTransactionUseCase(transactionRepo, userRepo, logger)

	consumer := mq.NewTransactionConsumer(rabbitConn, logger, useCase)

	go func() {
		_ = consumer.StartConsumer(
			cfg.AMQP.WorkerPoolSize,
			cfg.AMQP.Exchange,
			cfg.AMQP.Queue,
			cfg.AMQP.RoutingKey,
			cfg.AMQP.ConsumerTag,
		)

	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	logger.Info("received signal", slog.String("signal", s.String()))

}
