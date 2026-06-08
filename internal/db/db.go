package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func Connect() *Database {
	dsn := os.Getenv("DATABASE_URL")
	log.Println("Connecting to database with URL:")
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // ← disables prepared statements, fixes PgBouncer conflict
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal("failed to get SQL DB:", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("Database connected using URL")
	return &Database{
		DB: conn,
	}
}
