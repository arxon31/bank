package user

import (
	"context"
	"log/slog"

	"github.com/arxon31/bank/pkg/postgres"
)

type Repo struct {
	*postgres.Postgres
}

func NewRepo(postgres *postgres.Postgres) (*Repo, error) {
	return &Repo{Postgres: postgres}, nil
}

func (r *Repo) UpdateUserAmount(ctx context.Context, userID int64, amount int64) (err error) {
	const op = "user.UpdateUserAmount"
	logger := r.Logger.With(slog.String("op", op))

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	logger.Debug("updating user amount", slog.Int64("user_id", userID), slog.Int64("amount", amount))
	_, err = tx.ExecContext(ctx, "UPDATE users SET amount = $1, updated_at = NOW() WHERE id = $2", amount, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	logger.Debug("updated user amount", slog.Int64("user_id", userID), slog.Int64("amount", amount))
	return nil

}

func (r *Repo) GetUserAmount(ctx context.Context, userID int64) (amount int64, err error) {
	const op = "user.GetUserAmount"
	var currentAmount int64
	logger := r.Logger.With(slog.String("op", op))

	logger.Debug("getting user amount", slog.Int64("user_id", userID))

	row := r.DB.QueryRowContext(ctx, "SELECT amount FROM users WHERE id = $1", userID)
	err = row.Scan(&currentAmount)
	if err == nil {
		logger.Debug("got user amount", slog.Int64("amount", currentAmount))
	}
	return currentAmount, err
}
