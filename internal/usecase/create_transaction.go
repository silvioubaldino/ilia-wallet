package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/silvioubaldino/ilia-wallet/pkg/apperrors"
)

type createTransactionRepository interface {
	Create(ctx context.Context, t transaction.Transaction) (transaction.Transaction, error)
}

type CreateTransaction struct {
	repo createTransactionRepository
}

func NewCreateTransaction(repo createTransactionRepository) *CreateTransaction {
	return &CreateTransaction{repo: repo}
}

type CreateInput struct {
	UserID uuid.UUID
	Type   transaction.Type
	Amount int64
}

func (uc *CreateTransaction) Execute(ctx context.Context, input CreateInput) (transaction.Transaction, error) {
	if input.Type != transaction.TypeCredit && input.Type != transaction.TypeDebit {
		return transaction.Transaction{}, errors.Join(apperrors.ErrInvalidInput, errors.New("type must be CREDIT or DEBIT"))
	}
	if input.Amount <= 0 {
		return transaction.Transaction{}, errors.Join(apperrors.ErrInvalidInput, errors.New("amount must be greater than zero"))
	}

	t := transaction.Transaction{
		UserID: input.UserID,
		Type:   input.Type,
		Amount: input.Amount,
	}
	return uc.repo.Create(ctx, t)
}
