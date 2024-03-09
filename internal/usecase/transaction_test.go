package usecase

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/arxon31/bank/internal/entity"
	"github.com/golang/mock/gomock"
)

func TestTransaction(t *testing.T) {
	type fields struct {
		useCase  *MockTransaction
		userRepo *MockUserRepository

		transactionRepo *MockTransactionRepository
		logger          *slog.Logger
	}

	type args struct {
		ctx         context.Context
		transaction entity.Transaction
	}
	controller := gomock.NewController(t)

	transactionMock := NewMockTransaction(controller)
	transactionRepoMock := NewMockTransactionRepository(controller)
	userRepoMock := NewMockUserRepository(controller)
	logger := slog.Default()

	anyInput := gomock.Any()
	highBalance := int64(150)
	lowBalance := int64(10)
	someTransactionID := int64(1)

	var tests = []struct {
		name    string
		fields  fields
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(highBalance, nil)
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(lowBalance, nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(nil)
				f.transactionRepo.EXPECT().Store(anyInput, anyInput).Return(someTransactionID, nil)

			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: false,
		},
		{
			name: "Insufficient funds",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(lowBalance, nil)
			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: true,
		},
		{
			name: "Failed to get payer user amount",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(int64(0), errors.New("failed to get payer user amount"))
			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: true,
		},
		{
			name: "Failed to get payee user amount",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(highBalance, nil)
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(int64(0), errors.New("failed to get payer user amount"))
			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: true,
		},
		{
			name: "Failed to update payer user amount",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(highBalance, nil)
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(lowBalance, nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(errors.New("failed to update payer user amount"))
			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: true,
		},
		{
			name: "Failed to update payee user amount",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(highBalance, nil)
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(lowBalance, nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(errors.New("failed to update payee user amount"))
			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: true,
		},
		{
			name: "Failed to store transaction",
			fields: fields{
				useCase:         transactionMock,
				userRepo:        userRepoMock,
				transactionRepo: transactionRepoMock,
				logger:          logger,
			},
			prepare: func(f *fields) {
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(highBalance, nil)
				f.userRepo.EXPECT().GetUserAmount(anyInput, anyInput).Return(lowBalance, nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(nil)
				f.userRepo.EXPECT().UpdateUserAmount(anyInput, anyInput, anyInput).Return(nil)
				f.transactionRepo.EXPECT().Store(anyInput, anyInput).Return(int64(0), errors.New("failed to store transaction"))
			},
			args: args{
				ctx:         context.Background(),
				transaction: entity.Transaction{FromAccountID: 1, ToAccountID: 2, Amount: 100},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(&tt.fields)
			u := &TransactionUseCase{
				transactionRepo: tt.fields.transactionRepo,
				userRepo:        tt.fields.userRepo,
				logger:          tt.fields.logger,
			}
			_, err := u.MakeTransaction(tt.args.ctx, tt.args.transaction)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionUseCase.MakeTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
