package db

import (
	"log"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

func Migrate(db *Database) {
	err := db.DB.AutoMigrate(
		&models.Booking{},
		&models.Business{},
		&models.Service{},
		&models.User{},
		&models.BusinessHour{},
		&models.BlockedTime{},
		&models.Payment{},
		&models.BookingService{},
		&models.Ledger{},
		&models.OTP{},
	)

	if err != nil {
		log.Fatal("❌ migration failed:", err)
	}

	log.Println("✅ database migrated successfully")
}
