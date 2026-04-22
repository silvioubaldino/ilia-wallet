package handler_test

import (
	"context"

	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"github.com/stretchr/testify/mock"
)

type mockCreateUseCase struct {
	mock.Mock
}

func (m *mockCreateUseCase) Execute(_ context.Context, input usecase.CreateInput) (transaction.Transaction, error) {
	args := m.Called(input)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

type mockListUseCase struct {
	mock.Mock
}

func (m *mockListUseCase) Execute(_ context.Context, input usecase.ListInput) ([]transaction.Transaction, error) {
	args := m.Called(input)
	return args.Get(0).([]transaction.Transaction), args.Error(1)
}
