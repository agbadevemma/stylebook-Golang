package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Payment struct {
	Base
	BookingID  uuid.UUID `json:"booking_id" gorm:"type:uuid;not null;uniqueIndex"`
	BusinessID uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`

	// Paystack reference (VERY IMPORTANT for webhook idempotency)
	Reference string `json:"reference" gorm:"type:varchar(100);uniqueIndex;not null"`

	Name string `json:"name" gorm:"type:varchar(100);not null"`

	Email string `json:"email" gorm:"type:varchar(100);not null"`

	Amount int64 `json:"amount" gorm:"not null"` // total charged (kobo)

	Currency string `json:"currency" gorm:"type:varchar(10);default:'NGN'"`

	// breakdown (important for your platform model)
	ServiceAmount int64 `json:"service_amount" gorm:"not null"`
	PlatformFee   int64 `json:"platform_fee" gorm:"not null"`
	NetAmount     int64 `json:"net_amount" gorm:"not null"` // what barber receives

	Status utils.PaymentStatus `json:"status" gorm:"type:varchar(20);default:'pending';not null"`

	Provider string `json:"provider" gorm:"type:varchar(20);default:'paystack'"`

	// Paystack response metadata
	AuthorizationURL string `json:"authorization_url" gorm:"type:text"`
	AccessCode       string `json:"access_code" gorm:"type:varchar(100)"`

	// webhook safety (prevents double processing)
	IsVerified bool `json:"is_verified" gorm:"default:false"`

	// timestamps
	PaidAt *time.Time `json:"paid_at"`
}
