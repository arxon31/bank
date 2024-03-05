package app

import (
	"github.com/arxon31/bank/config"
	"github.com/arxon31/bank/internal/usecase"
	"github.com/arxon31/bank/internal/usecase/repo/transaction"
	"github.com/arxon31/bank/internal/usecase/repo/user"
	"github.com/arxon31/bank/pkg/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	const op = "app.Run"
	logger := slog.Default().With(slog.String("op", op))

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

	useCase := usecase.NewTransactionUseCase(transactionRepo, userRepo)

	//TODO:controller AMQP

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("received signal", slog.String("signal", s.String()))
	}

}
