package models

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Withdrawal struct {
	Base

	BusinessID uuid.UUID `gorm:"type:uuid;index;not null"`

	Amount int64
	Status utils.WithdrawalStatus `gorm:"type:varchar(20);not null;default:'pending'"` // pending | processing | success | failed

	RecipientCode string

	PaystackTransferCode string `gorm:"type:varchar(100);uniqueIndex"`
	TransferReference    string `gorm:"type:varchar(100);uniqueIndex;not null"`
}
