package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/utils"
	_ "github.com/go-sql-driver/mysql"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := &Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "3306"),
		User:     utils.GetEnv("DB_USER", "root"),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", "crypto_trading_db"),
	}

	// Debug: Print loaded configuration
	fmt.Printf("DEBUG: Loaded DB config - Host: %s, Port: %s, User: %s, DBName: %s\n",
		config.Host, config.Port, config.User, config.DBName)

	return config
}

// Connect establishes a connection to the MySQL database
func Connect(config *Config) (*sql.DB, error) {
	host := config.Host

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&timeout=10s&readTimeout=10s&writeTimeout=10s",
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.User,
		config.Password,
		host,
		config.Port,
		config.DBName,
	)

	fmt.Printf("DEBUG: Using DSN: %s\n", dsn)

	fmt.Printf("DEBUG: About to call sql.Open with DSN for host: %s\n", config.Host)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	fmt.Printf("DEBUG: sql.Open succeeded, about to ping database\n")

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
