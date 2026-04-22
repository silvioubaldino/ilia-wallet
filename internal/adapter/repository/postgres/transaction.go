package postgres

import (
	"context"

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
