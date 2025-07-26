package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Google   GoogleConfig   `mapstructure:"google"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	GinMode      string        `mapstructure:"gin_mode"`
}

type DatabaseConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Database     string        `mapstructure:"database"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	SSLMode      string        `mapstructure:"ssl_mode"`
	MaxConns     int           `mapstructure:"max_conns"`
	MinConns     int           `mapstructure:"min_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
	MaxIdleTime  time.Duration `mapstructure:"max_idle_time"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int           `mapstructure:"database"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	AccessExpiry   time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry  time.Duration `mapstructure:"refresh_expiry"`
	Issuer         string        `mapstructure:"issuer"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Set defaults
	setDefaults()
	
	// Set config file
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")
	}
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// Read environment variables
	viper.SetEnvPrefix("SMART_CALENDAR")
	viper.AutomaticEnv()
	
	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("server.idle_timeout", 120*time.Second)
	viper.SetDefault("server.gin_mode", "debug")
	
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.database", "smart_calendar")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_conns", 10)
	viper.SetDefault("database.min_conns", 2)
	viper.SetDefault("database.max_lifetime", time.Hour)
	viper.SetDefault("database.max_idle_time", 30*time.Minute)
	
	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", 5*time.Second)
	viper.SetDefault("redis.read_timeout", 3*time.Second)
	viper.SetDefault("redis.write_timeout", 3*time.Second)
	
	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.access_expiry", 15*time.Minute)
	viper.SetDefault("jwt.refresh_expiry", 24*time.Hour*7)
	viper.SetDefault("jwt.issuer", "smart-goal-calendar")
	
	// Google defaults
	viper.SetDefault("google.client_id", "")
	viper.SetDefault("google.client_secret", "")
	viper.SetDefault("google.redirect_url", "http://localhost:8080/auth/google/callback")
	
	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}