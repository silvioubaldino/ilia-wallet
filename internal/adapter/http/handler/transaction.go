package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/adapter/http/middleware"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"github.com/silvioubaldino/ilia-wallet/pkg/apperrors"
)

type createUseCase interface {
	Execute(ctx context.Context, input usecase.CreateInput) (transaction.Transaction, error)
}

type TransactionHandler struct {
	createUC createUseCase
}

func NewTransactionHandler(createUC createUseCase) *TransactionHandler {
	return &TransactionHandler{createUC: createUC}
}

type createRequest struct {
	UserID string           `json:"user_id" binding:"required"`
	Type   transaction.Type `json:"type"    binding:"required,oneof=CREDIT DEBIT"`
	Amount int64            `json:"amount"  binding:"required,gt=0"`
}

type transactionResponse struct {
	ID     string           `json:"id"`
	UserID string           `json:"user_id"`
	Type   transaction.Type `json:"type"`
	Amount int64            `json:"amount"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtUserID, _ := c.Get(middleware.UserIDKey)
	if req.UserID != jwtUserID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id mismatch"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	t, err := h.createUC.Execute(c.Request.Context(), usecase.CreateInput{
		UserID: userID,
		Type:   req.Type,
		Amount: req.Amount,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, transactionResponse{
		ID:     t.ID.String(),
		UserID: t.UserID.String(),
		Type:   t.Type,
		Amount: t.Amount,
	})
}
