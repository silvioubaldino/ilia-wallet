package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeCredit Type = "CREDIT"
	TypeDebit  Type = "DEBIT"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Type      Type      `gorm:"type:varchar(10);not null;check:type IN ('CREDIT','DEBIT')"`
	Amount    int64     `gorm:"not null;check:amount > 0"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}
