package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name string
		tx   Transaction
		err  error
	}{
		{
			name: "SameAccountIDs",
			tx: Transaction{
				FromAccountID: 1,
				ToAccountID:   1,
				Amount:        100,
			},
			err: ErrSameAccount,
		},
		{
			name: "NegativeAmount",
			tx: Transaction{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        -100,
			},
			err: ErrInvalidAmount,
		},
		{
			name: "InvalidFromAccountID",
			tx: Transaction{
				FromAccountID: 0,
				ToAccountID:   2,
				Amount:        100,
			},
			err: ErrFromID,
		},
		{
			name: "InvalidToAccountID",
			tx: Transaction{
				FromAccountID: 1,
				ToAccountID:   0,
				Amount:        100,
			},
			err: ErrToID,
		},
		{
			name: "ValidAccountIDs",
			tx: Transaction{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			assert.Equal(t, tt.err, err)
		})
	}
}
