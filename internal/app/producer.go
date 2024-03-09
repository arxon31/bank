package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/arxon31/bank/config/prod"
	"github.com/arxon31/bank/internal/controller/rabbtimq"
	"github.com/arxon31/bank/pkg/amqp"
	"github.com/arxon31/bank/pkg/logging"
)

func RunPublisher(cfg *prod.Config) {
	logger := logging.New(cfg.Env)

	ctx, cancel := context.WithCancel(context.Background())

	amqp, err := amqp.NewRabbitMQConn(cfg.AMQP.URL)
	if err != nil {
		logger.Error("amqp connection error", slog.String("url", cfg.AMQP.URL), slog.String("error", err.Error()))
		return
	}
	defer amqp.Close()
	logger.Info("amqp connection established", slog.String("url", cfg.AMQP.URL))

	publisher := rabbtimq.NewTransactionPublisher(amqp, logger, cfg.AMQP.Exchange, cfg.AMQP.RoutingKey)

	publisher.Start(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	logger.Info("received signal", slog.String("signal", s.String()))
	cancel()

}
