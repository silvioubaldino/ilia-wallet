package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestListTransactions_Execute(t *testing.T) {
	var (
		userID     = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		creditType = transaction.TypeCredit
		txList     = []transaction.Transaction{
			{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), UserID: userID, Type: transaction.TypeCredit, Amount: 100},
			{ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"), UserID: userID, Type: transaction.TypeDebit, Amount: 30},
		}
		creditList = []transaction.Transaction{txList[0]}
	)

	type input struct {
		userID uuid.UUID
		txType *transaction.Type
	}
	type mocks struct {
		repoList    []transaction.Transaction
		repoListErr error
	}
	type expected struct {
		output []transaction.Transaction
		err    error
	}

	tests := map[string]struct {
		input    input
		mocks    mocks
		expected expected
	}{
		"should return all transactions when no type filter": {
			input:    input{userID: userID, txType: nil},
			mocks:    mocks{repoList: txList, repoListErr: nil},
			expected: expected{output: txList, err: nil},
		},
		"should return only CREDIT transactions when type filter is CREDIT": {
			input:    input{userID: userID, txType: &creditType},
			mocks:    mocks{repoList: creditList, repoListErr: nil},
			expected: expected{output: creditList, err: nil},
		},
		"should return error when repo fails": {
			input:    input{userID: userID, txType: nil},
			mocks:    mocks{repoList: nil, repoListErr: assert.AnError},
			expected: expected{output: nil, err: assert.AnError},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &mockListTransactionRepository{}
			repo.On("List", tt.input.userID, tt.input.txType).Return(tt.mocks.repoList, tt.mocks.repoListErr)

			uc := usecase.NewListTransactions(repo)

			// Act
			got, err := uc.Execute(context.Background(), usecase.ListInput{
				UserID: tt.input.userID,
				Type:   tt.input.txType,
			})

			// Assert
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.output, got)
			repo.AssertExpectations(t)
		})
	}
}
