package user

import (
	"context"
	"github.com/arxon31/bank/pkg/postgres"
	"log/slog"
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
	_, err = tx.ExecContext(ctx, "UPDATE users SET amount = $1 WHERE id = $2", amount, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	logger.Info("updated user amount", slog.Int64("user_id", userID), slog.Int64("amount", amount))
	return nil

}

func (r *Repo) GetUserAmount(ctx context.Context, userID int64) (amount int64, err error) {
	const op = "user.GetUserAmount"
	logger := r.Logger.With(slog.String("op", op))

	logger.Debug("getting user amount", slog.Int64("user_id", userID))

	err = r.DB.QueryRowContext(ctx, "SELECT amount FROM users WHERE id = $1", userID).Scan(&amount)
	if err == nil {
		logger.Info("got user amount", slog.Int64("amount", amount))
	}
	return amount, err
}
