package usecase

import (
	"context"
	"errors"
	"github.com/arxon31/bank/internal/entity"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

const invalidTransactionID = -1

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
		return invalidTransactionID, err
	}

	if payerAmount < transaction.Amount {
		return invalidTransactionID, ErrInsufficientFunds
	}

	payeeAmount, err := tu.userRepo.GetUserAmount(ctx, transaction.ToAccountID)
	if err != nil {
		return invalidTransactionID, err
	}

	payerAmount -= transaction.Amount
	payeeAmount += transaction.Amount

	err = tu.userRepo.UpdateUserAmount(ctx, transaction.FromAccountID, payerAmount)
	if err != nil {
		return 0, err
	}

	err = tu.userRepo.UpdateUserAmount(ctx, transaction.ToAccountID, payeeAmount)
	if err != nil {
		return 0, err
	}

	return tu.transactionRepo.Store(ctx, transaction)

}
