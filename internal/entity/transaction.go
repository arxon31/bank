package entity

import "errors"

var (
	ErrSameAccount   = errors.New("same account")
	ErrInvalidAmount = errors.New("invalid amount")
	ErrFromID        = errors.New("invalid from account id")
	ErrToID          = errors.New("invalid to account id")
)

const InvalidTransactionID = -1

type Transaction struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (t Transaction) Validate() error {
	if t.FromAccountID == t.ToAccountID {
		return ErrSameAccount
	}
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	if t.FromAccountID < 1 {
		return ErrFromID
	}
	if t.ToAccountID < 1 {
		return ErrToID
	}
	return nil
}
