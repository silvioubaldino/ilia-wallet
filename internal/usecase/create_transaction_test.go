package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"github.com/silvioubaldino/ilia-wallet/pkg/apperrors"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction_Execute(t *testing.T) {
	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		stored = transaction.Transaction{
			ID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			UserID: userID,
			Type:   transaction.TypeCredit,
			Amount: 100,
		}
	)

	type input struct {
		userID uuid.UUID
		txType transaction.Type
		amount int64
	}
	type mocks struct {
		repoCreate         *transaction.Transaction
		repoCreateErr      error
		repoCreateCalled   bool
	}
	type expected struct {
		output transaction.Transaction
		err    error
	}

	tests := map[string]struct {
		input    input
		mocks    mocks
		expected expected
	}{
		"should return error when type is invalid": {
			input: input{userID: userID, txType: "INVALID", amount: 100},
			mocks: mocks{repoCreateCalled: false},
			expected: expected{
				output: transaction.Transaction{},
				err:    apperrors.ErrInvalidInput,
			},
		},
		"should return error when amount is zero": {
			input: input{userID: userID, txType: transaction.TypeCredit, amount: 0},
			mocks: mocks{repoCreateCalled: false},
			expected: expected{
				output: transaction.Transaction{},
				err:    apperrors.ErrInvalidInput,
			},
		},
		"should return error when amount is negative": {
			input: input{userID: userID, txType: transaction.TypeCredit, amount: -50},
			mocks: mocks{repoCreateCalled: false},
			expected: expected{
				output: transaction.Transaction{},
				err:    apperrors.ErrInvalidInput,
			},
		},
		"should return error when repo fails": {
			input: input{userID: userID, txType: transaction.TypeCredit, amount: 100},
			mocks: mocks{
				repoCreate:       &transaction.Transaction{},
				repoCreateErr:    assert.AnError,
				repoCreateCalled: true,
			},
			expected: expected{
				output: transaction.Transaction{},
				err:    assert.AnError,
			},
		},
		"should create transaction when input is valid": {
			input: input{userID: userID, txType: transaction.TypeCredit, amount: 100},
			mocks: mocks{
				repoCreate:       &stored,
				repoCreateErr:    nil,
				repoCreateCalled: true,
			},
			expected: expected{
				output: stored,
				err:    nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockCreateTransactionRepository{}
			if tt.mocks.repoCreateCalled {
				repo.On("Create", transaction.Transaction{
					UserID: tt.input.userID,
					Type:   tt.input.txType,
					Amount: tt.input.amount,
				}).Return(*tt.mocks.repoCreate, tt.mocks.repoCreateErr)
			}

			uc := usecase.NewCreateTransaction(repo)

			// Act
			got, err := uc.Execute(context.Background(), usecase.CreateInput{
				UserID: tt.input.userID,
				Type:   tt.input.txType,
				Amount: tt.input.amount,
			})

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
