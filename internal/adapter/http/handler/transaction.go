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

type listUseCase interface {
	Execute(ctx context.Context, input usecase.ListInput) ([]transaction.Transaction, error)
}

type balanceUseCase interface {
	Execute(ctx context.Context, input usecase.BalanceInput) (int64, error)
}

type TransactionHandler struct {
	createUC  createUseCase
	listUC    listUseCase
	balanceUC balanceUseCase
}

func NewTransactionHandler(createUC createUseCase, listUC listUseCase, balanceUC balanceUseCase) *TransactionHandler {
	return &TransactionHandler{createUC: createUC, listUC: listUC, balanceUC: balanceUC}
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

func (h *TransactionHandler) List(c *gin.Context) {
	jwtUserID, _ := c.Get(middleware.UserIDKey)
	userID, err := uuid.Parse(jwtUserID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
		return
	}

	var txType *transaction.Type
	if raw := c.Query("type"); raw != "" {
		t := transaction.Type(raw)
		if t != transaction.TypeCredit && t != transaction.TypeDebit {
			c.JSON(http.StatusBadRequest, gin.H{"error": "type must be CREDIT or DEBIT"})
			return
		}
		txType = &t
	}

	txs, err := h.listUC.Execute(c.Request.Context(), usecase.ListInput{
		UserID: userID,
		Type:   txType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]transactionResponse, len(txs))
	for i, t := range txs {
		resp[i] = transactionResponse{
			ID:     t.ID.String(),
			UserID: t.UserID.String(),
			Type:   t.Type,
			Amount: t.Amount,
		}
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) Balance(c *gin.Context) {
	jwtUserID, _ := c.Get(middleware.UserIDKey)
	userID, err := uuid.Parse(jwtUserID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
		return
	}

	amount, err := h.balanceUC.Execute(c.Request.Context(), usecase.BalanceInput{UserID: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"amount": amount})
}
