package db

import (
	"log"

	"github.com/stylebookgo/agbadevemma/internal/models"
)

func Migrate(db *Database) {
	err := db.DB.AutoMigrate(
		&models.Business{},
		&models.Service{},
		&models.Staff{},
		&models.Customer{},
		&models.Appointment{},
		&models.BusinessType{},
		&models.AppointmentService{},
		&models.WhatsAppSession{},
	)

	if err != nil {
		log.Fatal("❌ migration failed:", err)
	}

	log.Println("✅ database migrated successfully")
}
