package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	ID          int
	Name        string
	SQL         string
	AppliedAt   *time.Time
	Checksum    string
}

type Migrator struct {
	db         *pgxpool.Pool
	migrations []Migration
}

func NewMigrator(db *pgxpool.Pool) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []Migration{},
	}
}

// ensureMigrationsTable creates the migrations table if it doesn't exist
func (m *Migrator) ensureMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			checksum VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_schema_migrations_id ON schema_migrations(id);
	`
	
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	return nil
}

// LoadMigrationsFromDir loads migration files from a directory
func (m *Migrator) LoadMigrationsFromDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		migrationPath := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		// Parse migration ID from filename (e.g., "001_create_users.sql")
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			log.Printf("Warning: skipping migration file with invalid name format: %s", entry.Name())
			continue
		}

		var id int
		if _, err := fmt.Sscanf(parts[0], "%d", &id); err != nil {
			log.Printf("Warning: skipping migration file with invalid ID: %s", entry.Name())
			continue
		}

		migration := Migration{
			ID:       id,
			Name:     strings.TrimSuffix(entry.Name(), ".sql"),
			SQL:      string(content),
			Checksum: calculateChecksum(content),
		}

		m.migrations = append(m.migrations, migration)
	}

	// Sort migrations by ID
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	return nil
}

// GetAppliedMigrations returns list of already applied migrations
func (m *Migrator) GetAppliedMigrations() ([]Migration, error) {
	query := `
		SELECT id, name, checksum, applied_at 
		FROM schema_migrations 
		ORDER BY id
	`
	
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var applied []Migration
	for rows.Next() {
		var migration Migration
		var appliedAt time.Time
		
		err := rows.Scan(&migration.ID, &migration.Name, &migration.Checksum, &appliedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		
		migration.AppliedAt = &appliedAt
		applied = append(applied, migration)
	}

	return applied, nil
}

// Migrate runs all pending migrations
func (m *Migrator) Migrate() error {
	if err := m.ensureMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	appliedMap := make(map[int]Migration)
	for _, migration := range applied {
		appliedMap[migration.ID] = migration
	}

	for _, migration := range m.migrations {
		if appliedMigration, exists := appliedMap[migration.ID]; exists {
			// Check if migration content has changed
			if appliedMigration.Checksum != migration.Checksum {
				return fmt.Errorf("migration %d (%s) has been modified since it was applied", 
					migration.ID, migration.Name)
			}
			log.Printf("Migration %d (%s) already applied", migration.ID, migration.Name)
			continue
		}

		log.Printf("Applying migration %d (%s)...", migration.ID, migration.Name)
		
		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for migration %d: %w", migration.ID, err)
		}

		// Execute migration SQL
		_, err = tx.Exec(migration.SQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d (%s): %w", 
				migration.ID, migration.Name, err)
		}

		// Record migration as applied
		_, err = tx.Exec(
			"INSERT INTO schema_migrations (id, name, checksum) VALUES ($1, $2, $3)",
			migration.ID, migration.Name, migration.Checksum,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d as applied: %w", migration.ID, err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.ID, err)
		}

		log.Printf("Migration %d (%s) applied successfully", migration.ID, migration.Name)
	}

	return nil
}

// Status returns information about migration status
func (m *Migrator) Status() error {
	if err := m.ensureMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	appliedMap := make(map[int]Migration)
	for _, migration := range applied {
		appliedMap[migration.ID] = migration
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")
	
	for _, migration := range m.migrations {
		if appliedMigration, exists := appliedMap[migration.ID]; exists {
			status := "APPLIED"
			if appliedMigration.Checksum != migration.Checksum {
				status = "MODIFIED"
			}
			fmt.Printf("✓ %d %-40s %s (%s)\n", 
				migration.ID, migration.Name, status, appliedMigration.AppliedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("✗ %d %-40s PENDING\n", migration.ID, migration.Name)
		}
	}

	return nil
}

// calculateChecksum creates a simple checksum for migration content
func calculateChecksum(content []byte) string {
	// Simple hash of content length and first/last characters
	length := len(content)
	if length == 0 {
		return "empty"
	}
	
	first := content[0]
	last := content[length-1]
	
	return fmt.Sprintf("%d-%d-%d", length, first, last)
}