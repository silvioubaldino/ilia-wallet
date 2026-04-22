package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance_Execute(t *testing.T) {
	var (
		userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	)

	type input struct {
		userID uuid.UUID
	}
	type mocks struct {
		repoBalance    int64
		repoBalanceErr error
	}
	type expected struct {
		output int64
		err    error
	}

	tests := map[string]struct {
		input    input
		mocks    mocks
		expected expected
	}{
		"should return balance when repo succeeds": {
			input:    input{userID: userID},
			mocks:    mocks{repoBalance: 70, repoBalanceErr: nil},
			expected: expected{output: 70, err: nil},
		},
		"should return zero balance when no transactions": {
			input:    input{userID: userID},
			mocks:    mocks{repoBalance: 0, repoBalanceErr: nil},
			expected: expected{output: 0, err: nil},
		},
		"should return error when repo fails": {
			input:    input{userID: userID},
			mocks:    mocks{repoBalance: 0, repoBalanceErr: assert.AnError},
			expected: expected{output: 0, err: assert.AnError},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockBalanceRepository{}
			repo.On("Balance", tt.input.userID).Return(tt.mocks.repoBalance, tt.mocks.repoBalanceErr)

			uc := usecase.NewGetBalance(repo)

			// Act
			got, err := uc.Execute(context.Background(), usecase.BalanceInput{
				UserID: tt.input.userID,
			})

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
