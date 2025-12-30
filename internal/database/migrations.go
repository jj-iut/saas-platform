package database

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
		// Initial migration - create users table with roles
		migration := `
			-- Create enum type if it doesn't exist
			DO $$ BEGIN
				CREATE TYPE user_role AS ENUM ('user', 'admin', 'superadmin');
			EXCEPTION
				WHEN duplicate_object THEN null;
			END $$;
			
			-- Create users table if it doesn't exist
			CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				name VARCHAR(255),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);

			-- Add role column if it doesn't exist
			DO $$ BEGIN
				ALTER TABLE users ADD COLUMN role user_role DEFAULT 'user' NOT NULL;
			EXCEPTION
				WHEN duplicate_column THEN null;
			END $$;

			CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
			CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
		`

		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to apply migration 001_initial: %w", err)
		}

		// Record migration
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", "001_initial"); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
	}

	// Migration 002 - create restaurants table
	var count2 int
	err2 := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", "002_restaurants").Scan(&count2)
	if err2 != nil && err2 != sql.ErrNoRows {
		return fmt.Errorf("failed to check migration status: %w", err2)
	}

	if count2 == 0 {
		migration := `
			CREATE TABLE IF NOT EXISTS restaurants (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				address VARCHAR(500),
				phone VARCHAR(50),
				email VARCHAR(255),
				image_url VARCHAR(500),
				is_active BOOLEAN DEFAULT true,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);

			CREATE INDEX IF NOT EXISTS idx_restaurants_name ON restaurants(name);
			CREATE INDEX IF NOT EXISTS idx_restaurants_is_active ON restaurants(is_active);
		`

		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to apply migration 002_restaurants: %w", err)
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", "002_restaurants"); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
	}

	// Migration 003 - create default superadmin
	var count3 int
	err3 := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", "003_create_superadmin").Scan(&count3)
	if err3 != nil && err3 != sql.ErrNoRows {
		return fmt.Errorf("failed to check migration status: %w", err3)
	}

	if count3 == 0 {
		// Create default superadmin
		// Email: admin@saas-platform.com
		// Password: Admin123!
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash superadmin password: %w", err)
		}

		query := `
			INSERT INTO users (email, password_hash, role, name)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (email) DO NOTHING
		`

		if _, err := db.Exec(query, "admin@saas-platform.com", string(hashedPassword), "superadmin", "Super Admin"); err != nil {
			return fmt.Errorf("failed to create superadmin: %w", err)
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", "003_create_superadmin"); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}
	}

	return nil
}
