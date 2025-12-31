package config

import (
    "os"
    "path/filepath"
)

type Config struct {
    DatabasePath string
    ServerPort   string
}

func LoadConfig() *Config {
    // Get executable directory
    exePath, err := os.Executable()
    if err != nil {
        exePath = "."
    }
    exeDir := filepath.Dir(exePath)

    // Default database path
    dbPath := filepath.Join(exeDir, "database", "posflow.db")
    if dbPathEnv := os.Getenv("DB_PATH"); dbPathEnv != "" {
        dbPath = dbPathEnv
    }

    // Default server port
    port := "8000"
    if portEnv := os.Getenv("PORT"); portEnv != "" {
        port = portEnv
    }

    return &Config{
        DatabasePath: dbPath,
        ServerPort:   port,
    }
}