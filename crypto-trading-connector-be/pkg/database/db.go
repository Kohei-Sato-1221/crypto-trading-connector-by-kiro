package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-trading-connector/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DBType   string
	DSN      string
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := &Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "3306"),
		User:     utils.GetEnv("DB_USER", "root"),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", "crypto_trading_db"),
		DBType:   utils.GetEnv("DB_TYPE", "mysql"),
		DSN:      utils.GetEnv("DB_DSN", ""),
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

// ConnectPostgres establishes a connection to a PostgreSQL database
func ConnectPostgres(config *Config) (*sql.DB, error) {
	var dsn string
	if config.DSN != "" {
		dsn = config.DSN
	} else {
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require",
			config.User, config.Password, config.Host, config.Port, config.DBName)
	}

	fmt.Printf("DEBUG: Connecting to PostgreSQL at host: %s\n", config.Host)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	return db, nil
}

// ConnectByType establishes a database connection based on the configured DB_TYPE
func ConnectByType(config *Config) (*sql.DB, error) {
	switch config.DBType {
	case "postgres":
		return ConnectPostgres(config)
	case "mysql":
		return Connect(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.DBType)
	}
}
