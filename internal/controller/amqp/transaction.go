package amqp

import (
	"github.com/arxon31/bank/internal/usecase"
)

type transactionRoutes struct {
	transactionUseCase usecase.Transaction
}

func New(transactionUseCase usecase.Transaction) *transactionRoutes {

}
