package usecase

import (
	"context"
	"errors"
	"github.com/arxon31/bank/internal/entity"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type TransactionUseCase struct {
	transactionRepo TransactionRepository
	userRepo        UserRepository
}

func NewTransactionUseCase(tr TransactionRepository, ur UserRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: tr,
		userRepo:        ur,
	}
}

func (tu *TransactionUseCase) MakeTransaction(ctx context.Context, transaction entity.Transaction) (transactionID int64, err error) {
	payerAmount, err := tu.userRepo.GetUserAmount(ctx, transaction.FromAccountID)
	if err != nil {
		return entity.InvalidTransactionID, err
	}

	if payerAmount < transaction.Amount {
		return entity.InvalidTransactionID, ErrInsufficientFunds
	}

	payeeAmount, err := tu.userRepo.GetUserAmount(ctx, transaction.ToAccountID)
	if err != nil {
		return entity.InvalidTransactionID, err
	}

	payerAmount -= transaction.Amount
	payeeAmount += transaction.Amount

	err = tu.userRepo.UpdateUserAmount(ctx, transaction.FromAccountID, payerAmount)
	if err != nil {
		return entity.InvalidTransactionID, err
	}

	err = tu.userRepo.UpdateUserAmount(ctx, transaction.ToAccountID, payeeAmount)
	if err != nil {
		return entity.InvalidTransactionID, err
	}

	return tu.transactionRepo.Store(ctx, transaction)

}
