package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
)

type listTransactionRepository interface {
	List(ctx context.Context, userID uuid.UUID, txType *transaction.Type) ([]transaction.Transaction, error)
}

type ListTransactions struct {
	repo listTransactionRepository
}

func NewListTransactions(repo listTransactionRepository) *ListTransactions {
	return &ListTransactions{repo: repo}
}

type ListInput struct {
	UserID uuid.UUID
	Type   *transaction.Type
}

func (uc *ListTransactions) Execute(ctx context.Context, input ListInput) ([]transaction.Transaction, error) {
	return uc.repo.List(ctx, input.UserID, input.Type)
}
