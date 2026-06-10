package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Ledger struct {
	ID         uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	BusinessID uuid.UUID        `gorm:"type:uuid;index;not null"`
	Amount     int64            `json:"platform_fee" gorm:"not null"` 
	Type       utils.LedgerType `gorm:"type:varchar(20);not null"` // positive = credit, negative = debit
	Reference  string           `gorm:"type:varchar(100);index"` // payment/ref/withdrawal id
	Reason     string           `gorm:"type:varchar(255)"`
	CreatedAt  time.Time
}
