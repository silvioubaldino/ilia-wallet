package usecase

import (
	"context"

	"github.com/google/uuid"
)

type balanceRepository interface {
	Balance(ctx context.Context, userID uuid.UUID) (int64, error)
}

type GetBalance struct {
	repo balanceRepository
}

func NewGetBalance(repo balanceRepository) *GetBalance {
	return &GetBalance{repo: repo}
}

type BalanceInput struct {
	UserID uuid.UUID
}

func (uc *GetBalance) Execute(ctx context.Context, input BalanceInput) (int64, error) {
	return uc.repo.Balance(ctx, input.UserID)
}
