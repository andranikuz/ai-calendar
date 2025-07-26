package main

import (
	"log"

	"github.com/andranikuz/smart-goal-calendar/config"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/migrations"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/postgres"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.NewDatabase(postgres.Config{
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
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	migrator := migrations.NewMigrator(db.Pool)
	if err := migrator.LoadMigrationsFromDir("migrations"); err != nil {
		log.Fatalf("failed to load migrations: %v", err)
	}

	if err := migrator.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
}
