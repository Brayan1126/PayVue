package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabasePath       string
	Environment        string
	LogLevel           string
	CORSAllowedOrigins string
	ServerTimeout      int
}

func init() {
	// Cargar el archivo .env
	envPath := findEnvFile()

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Could not load .env file from %s: %v. Using system environment variables or defaults.", envPath, err)
	} else {
		absPath, _ := filepath.Abs(envPath)
		log.Printf("Loaded configuration from: %s", absPath)
	}
}

func findEnvFile() string {
	// Buscar el archivo .env en diferentes ubicaciones
	possiblePaths := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Si no se encuentra, usar la ruta por defecto
	return ".env"
}

func Get() Config {
	port := getEnv("PORT", "8080")
	databasePath := getEnv("DATABASE_PATH", "./payvue.db")
	environment := getEnv("ENVIRONMENT", "development")
	logLevel := getEnv("LOG_LEVEL", "info")
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "*")

	timeoutStr := getEnv("SERVER_TIMEOUT", "60")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 60
	}

	return Config{
		Port:               port,
		DatabasePath:       databasePath,
		Environment:        environment,
		LogLevel:           logLevel,
		CORSAllowedOrigins: corsOrigins,
		ServerTimeout:      timeout,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
