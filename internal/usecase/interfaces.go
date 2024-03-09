package usecase

import (
	"context"

	"github.com/arxon31/bank/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase

type (
	Transaction interface {
		MakeTransaction(ctx context.Context, transaction entity.Transaction) (transactionID int64, err error)
	}

	TransactionRepository interface {
		Store(ctx context.Context, transaction entity.Transaction) (transactionID int64, err error)
	}

	UserRepository interface {
		GetUserAmount(ctx context.Context, userID int64) (amount int64, err error)
		UpdateUserAmount(ctx context.Context, userID int64, amount int64) (err error)
	}
)
