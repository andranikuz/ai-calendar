package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andranikuz/smart-goal-calendar/config"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/migrations"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/postgres"
)

func main() {
	var (
		configPath = flag.String("config", "", "path to config file")
		action     = flag.String("action", "migrate", "action to perform: migrate, status, or create")
		name       = flag.String("name", "", "name for new migration (only for create action)")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbConfig := postgres.Config{
		Host:        cfg.Database.Host,
		Port:        cfg.Database.Port,
		Database:    cfg.Database.Database,
		User:        cfg.Database.User,
		Password:    cfg.Database.Password,
		SSLMode:     cfg.Database.SSLMode,
		MaxConns:    cfg.Database.MaxConns,
		MinConns:    cfg.Database.MinConns,
		MaxLifetime: cfg.Database.MaxLifetime,
		MaxIdleTime: cfg.Database.MaxIdleTime,
	}

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Health(ctx); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	fmt.Println("Database connected successfully")

	// Initialize migrator
	migrator := migrations.NewMigrator(db.Pool)

	// Load migrations from directory
	if err := migrator.LoadMigrationsFromDir("migrations"); err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	// Perform action
	switch *action {
	case "migrate":
		fmt.Println("Running database migrations...")
		if err := migrator.Migrate(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully!")

	case "status":
		fmt.Println("Checking migration status...")
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

	case "create":
		if *name == "" {
			log.Fatal("Migration name is required for create action")
		}
		if err := createMigration(*name); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		fmt.Printf("Migration created successfully: %s\n", *name)

	default:
		log.Fatalf("Unknown action: %s. Use migrate, status, or create", *action)
	}
}

func createMigration(name string) error {
	// Find the next migration ID
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	maxID := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		var id int
		if n, _ := fmt.Sscanf(entry.Name(), "%d_", &id); n == 1 && id > maxID {
			maxID = id
		}
	}

	nextID := maxID + 1
	filename := fmt.Sprintf("migrations/%03d_%s.sql", nextID, name)

	content := fmt.Sprintf(`-- Migration %03d: %s
-- Add your SQL here

`, nextID, name)

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	return nil
}
