package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-wallet/internal/domain"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t transaction.Transaction) (transaction.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(&t).Error; err != nil {
		return transaction.Transaction{}, err
	}
	return t, nil
}

func (r *TransactionRepository) List(ctx context.Context, userID uuid.UUID, txType *transaction.Type) ([]transaction.Transaction, error) {
	var txs []transaction.Transaction
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if txType != nil {
		q = q.Where("type = ?", *txType)
	}
	if err := q.Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
