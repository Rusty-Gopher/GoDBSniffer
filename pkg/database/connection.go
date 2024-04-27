package database

import (
	"GoDBSniffer/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig holds the database configuration details
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// ConnectDatabase establishes a connection to the database using the provided configuration.
func ConnectDatabase(config config.DBConfig) (*sql.DB, error) {
	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Database)

	// Open a new connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}
