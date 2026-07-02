package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Ledger struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	BusinessID  uuid.UUID        `gorm:"type:uuid;index;not null"`
	NetAmount   int64            `json:"netAmount" gorm:"not null;default:0"`                        // amount after deducting platform fee
	PlatformFee int64            `json:"platformFee" gorm:"not null;default:0"`                      // platform fee deducted from the total amount
	TotalAmount int64            `json:"totalAmount" gorm:"not null;default:0"`                      // total amount before deducting platform fee
	Type        utils.LedgerType `gorm:"type:varchar(20);not null;uniqueIndex:idx_ledger_ref_type"`  // positive = credit, negative = debit
	Balance     int64            `json:"balance" gorm:"not null;default:0"`                          // balance after this transaction
	Reference   string           `gorm:"type:varchar(100);not null;uniqueIndex:idx_ledger_ref_type"` // payment/ref/withdrawal id
	Reason      string           `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
}
