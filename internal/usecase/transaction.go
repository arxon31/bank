package usecase

import (
	"context"
	"errors"
	"github.com/arxon31/bank/internal/entity"
	"log/slog"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type TransactionUseCase struct {
	transactionRepo TransactionRepository
	userRepo        UserRepository
	logger          *slog.Logger
}

func NewTransactionUseCase(tr TransactionRepository, ur UserRepository, logger *slog.Logger) TransactionUseCase {
	return TransactionUseCase{
		transactionRepo: tr,
		userRepo:        ur,
		logger:          logger,
	}
}

func (tu *TransactionUseCase) MakeTransaction(ctx context.Context, transaction entity.Transaction) (transactionID int64, err error) {
	payerAmount, err := tu.userRepo.GetUserAmount(ctx, transaction.FromAccountID)
	if err != nil {
		tu.logger.Error("failed to get payer user amount", slog.String("error", err.Error()))
		return entity.InvalidTransactionID, err
	}
	tu.logger.Debug("got user amount", slog.Int64("user_id", transaction.FromAccountID), slog.Int64("amount", payerAmount))

	if payerAmount < transaction.Amount {
		tu.logger.Debug("insufficient funds", slog.Int64("user_id", transaction.FromAccountID), slog.Int64("amount", payerAmount), slog.Int64("transaction_amount", transaction.Amount))
		return entity.InvalidTransactionID, ErrInsufficientFunds
	}

	payeeAmount, err := tu.userRepo.GetUserAmount(ctx, transaction.ToAccountID)
	if err != nil {
		tu.logger.Error("failed to get payee user amount", slog.String("error", err.Error()))
		return entity.InvalidTransactionID, err
	}

	payerAmount -= transaction.Amount
	payeeAmount += transaction.Amount

	err = tu.userRepo.UpdateUserAmount(ctx, transaction.FromAccountID, payerAmount)
	if err != nil {
		tu.logger.Error("failed to update payer user amount", slog.String("error", err.Error()))
		return entity.InvalidTransactionID, err
	}

	err = tu.userRepo.UpdateUserAmount(ctx, transaction.ToAccountID, payeeAmount)
	if err != nil {
		tu.logger.Error("failed to update payee user amount", slog.String("error", err.Error()))
		return entity.InvalidTransactionID, err
	}

	tu.logger.Info("successful transaction")

	return tu.transactionRepo.Store(ctx, transaction)

}
