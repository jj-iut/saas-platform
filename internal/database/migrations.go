package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check if migrations have already been applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", "001_initial").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if count == 0 {
		// Initial migration - create users table as an example
		migration := `
			CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);

			CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		`

		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to apply migration 001_initial: %w", err)
		}

		// Record migration
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", "001_initial"); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
	}

	return nil
}
