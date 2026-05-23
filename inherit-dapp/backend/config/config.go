package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	// Server
	ServerPort int    `json:"server_port"`
	ServerHost string `json:"server_host"`

	// Database
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBSSLMode  string `json:"db_ssl_mode"`

	// Blockchain
	ChainRPCEndpoint  string `json:"chain_rpc_endpoint"`
	ChainGRPCEndpoint string `json:"chain_grpc_endpoint"`
	ChainID           string `json:"chain_id"`
	KeyringBackend    string `json:"keyring_backend"`
	NodeHome          string `json:"node_home"`

	// IPFS
	IPFSAPIEndpoint string `json:"ipfs_api_endpoint"`

	// Monitor
	HeartbeatCheckInterval time.Duration `json:"heartbeat_check_interval"`
	MonitorEnabled         bool          `json:"monitor_enabled"`

	// Security
	JWTSecret string `json:"jwt_secret"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnvInt("SERVER_PORT", 8080),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "legacy"),
		DBPassword: getEnv("DB_PASSWORD", "legacy_password"),
		DBName:     getEnv("DB_NAME", "legacy_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		ChainRPCEndpoint:  getEnv("CHAIN_RPC_ENDPOINT", "http://localhost:26657"),
		ChainGRPCEndpoint: getEnv("CHAIN_GRPC_ENDPOINT", "localhost:9090"),
		ChainID:           getEnv("CHAIN_ID", "legacy-testnet-1"),
		KeyringBackend:    getEnv("KEYRING_BACKEND", "test"),
		NodeHome:          getEnv("NODE_HOME", "~/.legacyd"),

		IPFSAPIEndpoint: getEnv("IPFS_API_ENDPOINT", "http://localhost:5001/api/v0"),

		HeartbeatCheckInterval: time.Duration(getEnvInt("HEARTBEAT_CHECK_INTERVAL", 60)) * time.Second,
		MonitorEnabled:         getEnv("MONITOR_ENABLED", "true") == "true",

		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + strconv.Itoa(c.DBPort) +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
