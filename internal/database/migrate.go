package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// RunMigrations executes the migration files to setup the database schema
func RunMigrations(db *sql.DB) error {
	// Read the migration file
	migrationPath := filepath.Join("migrations", "001_init_schema.up.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute the migration
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Database schema created successfully")
	return nil
}