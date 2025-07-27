package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/andranikuz/smart-goal-calendar/config"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/auth"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/google"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/migrations"
	"github.com/andranikuz/smart-goal-calendar/internal/adapters/postgres"
	appHandlers "github.com/andranikuz/smart-goal-calendar/internal/application/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/services"
	httpHandlers "github.com/andranikuz/smart-goal-calendar/internal/ports/http/handlers"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/middleware"
	"github.com/andranikuz/smart-goal-calendar/internal/ports/http/routes"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	setupLogger(cfg.Logging)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

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
		zlog.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Health(ctx); err != nil {
		zlog.Fatal().Err(err).Msg("Database health check failed")
	}

	zlog.Info().Msg("Database connected successfully")

	// Run database migrations
	migrator := migrations.NewMigrator(db.Pool)
	if err := migrator.LoadMigrationsFromDir("migrations"); err != nil {
		zlog.Fatal().Err(err).Msg("Failed to load migrations")
	}
	
	if err := migrator.Migrate(); err != nil {
		zlog.Fatal().Err(err).Msg("Failed to run migrations")
	}
	
	zlog.Info().Msg("Database migrations completed successfully")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db.Pool)
	goalRepo := postgres.NewGoalRepository(db.Pool)
	taskRepo := postgres.NewTaskRepository(db.Pool)
	milestoneRepo := postgres.NewMilestoneRepository(db.Pool)
	eventRepo := postgres.NewEventRepository(db.Pool)
	moodRepo := postgres.NewMoodRepository(db.Pool)
	googleIntegrationRepo := postgres.NewGoogleIntegrationRepository(db.Pool)
	googleCalendarSyncRepo := postgres.NewGoogleCalendarSyncRepository(db.Pool)

	// Initialize services
	userService := services.NewUserService()
	goalService := services.NewGoalService()
	eventService := services.NewEventService()
	moodService := services.NewMoodService()

	// Initialize application handlers
	userHandler := appHandlers.NewUserHandler(userRepo)
	goalHandler := appHandlers.NewGoalHandler(goalRepo, taskRepo, milestoneRepo, goalService)
	eventHandler := appHandlers.NewEventHandler(eventRepo, goalRepo, eventService)
	moodHandler := appHandlers.NewMoodHandler(moodRepo, moodService)

	// Initialize JWT service
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
		cfg.JWT.Issuer,
	)

	// Initialize Google services
	oauth2Service := google.NewOAuth2Service(
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
	)
	calendarService := google.NewCalendarService(oauth2Service)

	// Initialize HTTP handlers
	userHTTPHandler := httpHandlers.NewUserHTTPHandler(userHandler, userService, jwtService)
	goalHTTPHandler := httpHandlers.NewGoalHTTPHandler(goalHandler)
	eventHTTPHandler := httpHandlers.NewEventHTTPHandler(eventHandler)
	moodHTTPHandler := httpHandlers.NewMoodHTTPHandler(moodHandler)
	googleAuthHandler := httpHandlers.NewGoogleAuthHandler(
		oauth2Service,
		calendarService,
		googleIntegrationRepo,
		googleCalendarSyncRepo,
	)
	googleCalendarSyncHandler := httpHandlers.NewGoogleCalendarSyncHandler(
		oauth2Service,
		calendarService,
		googleIntegrationRepo,
		googleCalendarSyncRepo,
		eventRepo,
	)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup HTTP router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "smart-goal-calendar",
			"version": "0.1.0",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Setup user routes
		routes.SetupUserRoutes(v1, userHTTPHandler, authMiddleware)

		// Setup goal routes
		routes.SetupGoalRoutes(v1, goalHTTPHandler, authMiddleware)

		// Setup event routes
		routes.SetupEventRoutes(v1, eventHTTPHandler, authMiddleware)

		// Setup mood routes
		routes.SetupMoodRoutes(v1, moodHTTPHandler, authMiddleware)

		// Setup Google authentication routes
		routes.SetupGoogleAuthRoutes(v1, googleAuthHandler, authMiddleware)

		// Setup Google calendar sync routes
		routes.SetupGoogleCalendarSyncRoutes(v1, googleCalendarSyncHandler, authMiddleware)
	}

	// Setup HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		zlog.Info().
			Str("host", cfg.Server.Host).
			Int("port", cfg.Server.Port).
			Msg("Starting HTTP server")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zlog.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	zlog.Info().Msg("Server exited")
}

func setupLogger(cfg config.LoggingConfig) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set log format
	if cfg.Format == "console" {
		zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Add timestamp
	zerolog.TimeFieldFormat = time.RFC3339
}
