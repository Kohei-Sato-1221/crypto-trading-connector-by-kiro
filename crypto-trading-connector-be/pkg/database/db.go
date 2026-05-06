package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/crypto-trading-connector/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	MysqlURL    string // MySQL接続URL (例: user:pass@tcp(host:port)/dbname?parseTime=true&loc=Local)
	PostgresURL string // PostgreSQL接続URL (例: postgresql://user:pass@host:port/dbname?sslmode=require)
	DBType      string // 自動判定結果: "mysql" or "postgres"
}

// LoadConfigFromEnv loads database configuration from environment variables.
// POSTGRES_URL が設定されていればPostgreSQLを優先、MYSQL_URLのみならMySQL。
// 両方空の場合はpanicする。
func LoadConfigFromEnv() *Config {
	mysqlURL := utils.GetEnv("MYSQL_URL", "")
	postgresURL := utils.GetEnv("POSTGRES_URL", "")

	if mysqlURL == "" && postgresURL == "" {
		log.Fatal("FATAL: Neither MYSQL_URL nor POSTGRES_URL is set. At least one database URL must be configured.")
	}

	// PostgreSQL優先
	dbType := "mysql"
	if postgresURL != "" {
		dbType = "postgres"
	}

	log.Printf("DEBUG: DB config - DBType: %s, MYSQL_URL set: %v, POSTGRES_URL set: %v",
		dbType, mysqlURL != "", postgresURL != "")

	return &Config{
		MysqlURL:    mysqlURL,
		PostgresURL: postgresURL,
		DBType:      dbType,
	}
}

// Connect establishes a database connection based on which URL is configured.
// PostgreSQL URL が設定されていればPostgreSQLを優先する。
func Connect(config *Config) (*sql.DB, error) {
	switch config.DBType {
	case "postgres":
		return connectPostgres(config.PostgresURL)
	default:
		return connectMySQL(config.MysqlURL)
	}
}

// connectMySQL establishes a connection to a MySQL database
func connectMySQL(dsn string) (*sql.DB, error) {
	log.Printf("DEBUG: Connecting to MySQL")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	return db, nil
}

// connectPostgres establishes a connection to a PostgreSQL database
func connectPostgres(dsn string) (*sql.DB, error) {
	log.Printf("DEBUG: Connecting to PostgreSQL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}

	return db, nil
}
