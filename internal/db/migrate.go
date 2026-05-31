package db

import (
	"log"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

func Migrate(db *Database) {
	err := db.DB.AutoMigrate(
		&models.Booking{},
		&models.Client{},
		&models.Business{},
		&models.Service{},
		&models.User{},
		&models.Availability{},
		&models.BlockedTime{},
	)

	if err != nil {
		log.Fatal("❌ migration failed:", err)
	}

	log.Println("✅ database migrated successfully")
}
