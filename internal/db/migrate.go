package db

import (
	"log"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

func Migrate(db *Database) {
	err := db.DB.AutoMigrate(
		&models.Barber{},
		&models.Booking{},
		&models.Client{},
		&models.User{},
	)

	if err != nil {
		log.Fatal("❌ migration failed:", err)
	}

	log.Println("✅ database migrated successfully")
}
