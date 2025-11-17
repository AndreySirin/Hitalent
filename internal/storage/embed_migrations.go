package storage

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MigrateUP(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	err = goose.Up(sqlDB, "migrations")
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}
