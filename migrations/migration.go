package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
)

//go:embed postgres/*.sql
var embedMigrations embed.FS

func Migration(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(db, "postgres"); err != nil {
		return fmt.Errorf("миграция не прошла: %w", err)
	}
	return nil
}
