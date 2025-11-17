package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
)

func New(user, password, address, dbname string) (*gorm.DB, error) {

	dsn := (&url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(user, password),
		Host:   address,
		Path:   dbname,
	}).String()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failid connect: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}
