package models

import "github.com/google/uuid"

type Business struct {
	Base
	Name     string    `json:"name" gorm:"type:varchar(255);not null"`
	Timezone string    `json:"timezone" gorm:"type:varchar(50);default:'Africa/Lagos'"`
	OwnerID  uuid.UUID `json:"ownerId" gorm:"type:varchar(255);not null;index"`
	Owner    User      `json:"owner" gorm:"foreignKey:OwnerID"`

	AccountNumber     string `json:"accountNumber" gorm:"type:varchar(20);default:''"`
	BankName          string `json:"bankName" gorm:"type:varchar(100);default:''"`
	BankCode          string `json:"bankCode" gorm:"type:varchar(10);default:''"`
	AccountName       string `json:"accountName" gorm:"type:varchar(900);default:''"`
	TransferRecipient string `json:"transferRecipient" gorm:"type:varchar(100);default:''"`
	
	WithdrawalPIN     string `json:"withdrawalPIN" gorm:"type:varchar(255);default:''"` // hashed PIN for withdrawal verification
}
