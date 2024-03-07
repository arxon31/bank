package transaction

import (
	"context"
	"github.com/arxon31/bank/internal/entity"
	"github.com/arxon31/bank/pkg/postgres"
	"log/slog"
)

type Repo struct {
	*postgres.Postgres
}

func NewRepo(postgres *postgres.Postgres) (*Repo, error) {

	return &Repo{Postgres: postgres}, nil
}

func (r *Repo) Store(ctx context.Context, transaction entity.Transaction) (transactionID int64, err error) {
	const op = "transaction.Store"

	logger := r.Logger.With(slog.String("op", op))

	var id int64

	tx, err := r.DB.Begin()
	if err != nil {
		return entity.InvalidTransactionID, err
	}
	logger.Debug("inserting transaction", slog.Int64("from_account_id", transaction.FromAccountID), slog.Int64("to_account_id", transaction.ToAccountID), slog.Int64("amount", transaction.Amount))

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (from_id, to_id, amount) VALUES ($1, $2, $3)", transaction.FromAccountID, transaction.ToAccountID, transaction.Amount)
	if err != nil {
		tx.Rollback()
		return entity.InvalidTransactionID, err
	}

	logger.Debug("inserted transaction", slog.Int64("id", id))

	err = tx.Commit()
	if err != nil {
		return entity.InvalidTransactionID, err
	}

	logger.Info("inserted transaction", slog.Int64("id", id))

	return id, nil

}
