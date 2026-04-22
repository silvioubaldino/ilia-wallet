package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/adapter/http/middleware"
	"github.com/silvioubaldino/ilia-wallet/internal/adapter/http/handler"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var (
		userID  = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		transID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
		stored  = transaction.Transaction{
			ID:     transID,
			UserID: userID,
			Type:   transaction.TypeCredit,
			Amount: 100,
		}
	)

	type inputBody struct {
		UserID string `json:"user_id"`
		Type   string `json:"type"`
		Amount int64  `json:"amount"`
	}
	type mocks struct {
		ucInput    *usecase.CreateInput
		ucOutput   *transaction.Transaction
		ucErr      error
		ucCalled   bool
	}
	type expected struct {
		statusCode int
	}

	tests := map[string]struct {
		inputBody      inputBody
		jwtUserID      string
		mocks          mocks
		expected       expected
	}{
		"should return 400 when body is missing required fields": {
			inputBody: inputBody{},
			jwtUserID: userID.String(),
			mocks:     mocks{ucCalled: false},
			expected:  expected{statusCode: http.StatusBadRequest},
		},
		"should return 400 when type is invalid": {
			inputBody: inputBody{UserID: userID.String(), Type: "INVALID", Amount: 100},
			jwtUserID: userID.String(),
			mocks:     mocks{ucCalled: false},
			expected:  expected{statusCode: http.StatusBadRequest},
		},
		"should return 401 when user_id does not match jwt claim": {
			inputBody: inputBody{UserID: uuid.New().String(), Type: "CREDIT", Amount: 100},
			jwtUserID: userID.String(),
			mocks:     mocks{ucCalled: false},
			expected:  expected{statusCode: http.StatusUnauthorized},
		},
		"should return 500 when usecase fails": {
			inputBody: inputBody{UserID: userID.String(), Type: "CREDIT", Amount: 100},
			jwtUserID: userID.String(),
			mocks: mocks{
				ucInput:  &usecase.CreateInput{UserID: userID, Type: transaction.TypeCredit, Amount: 100},
				ucOutput: &transaction.Transaction{},
				ucErr:    assert.AnError,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusInternalServerError},
		},
		"should return 200 when transaction is created": {
			inputBody: inputBody{UserID: userID.String(), Type: "CREDIT", Amount: 100},
			jwtUserID: userID.String(),
			mocks: mocks{
				ucInput:  &usecase.CreateInput{UserID: userID, Type: transaction.TypeCredit, Amount: 100},
				ucOutput: &stored,
				ucErr:    nil,
				ucCalled: true,
			},
			expected: expected{statusCode: http.StatusOK},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			uc := &mockCreateUseCase{}
			if tt.mocks.ucCalled {
				uc.On("Execute", *tt.mocks.ucInput).Return(*tt.mocks.ucOutput, tt.mocks.ucErr)
			}

			h := handler.NewTransactionHandler(uc)

			router := gin.New()
			router.POST("/transactions", func(c *gin.Context) {
				c.Set(middleware.UserIDKey, tt.jwtUserID)
				h.Create(c)
			})

			body, _ := json.Marshal(tt.inputBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expected.statusCode, w.Code)
			uc.AssertExpectations(t)
		})
	}
}
