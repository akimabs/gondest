package config

import (
	"fmt"
	"os"

	"gorm.io/driver/{{ .DatabaseDriver }}"
	"gorm.io/gorm"
)

func NewGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

    db, err := gorm.Open({{ .DatabaseDriver }}.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}