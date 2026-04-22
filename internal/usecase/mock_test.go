package usecase_test

import (
	"context"

	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/stretchr/testify/mock"
)

type mockCreateTransactionRepository struct {
	mock.Mock
}

func (m *mockCreateTransactionRepository) Create(_ context.Context, t transaction.Transaction) (transaction.Transaction, error) {
	args := m.Called(t)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}
