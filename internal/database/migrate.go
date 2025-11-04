package database


import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// RunMigrations executes the complete schema
func RunMigrations(db *sql.DB) error {
	// Read the schema file
	schemaPath := filepath.Join("database", "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute the schema
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database schema created successfully")
	return nil
}