// config/config.go
package config

// DBConfig holds the database configuration details.
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// dbConfig is a package-level variable that holds the active DB configuration.
var dbConfig DBConfig

// SetDBConfig sets the active DB configuration.
func SetDBConfig(cfg DBConfig) {
	dbConfig = cfg
}

// GetDBConfig retrieves the active DB configuration.
func GetDBConfig() DBConfig {
	return dbConfig
}
