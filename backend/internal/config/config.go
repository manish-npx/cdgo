package config

import (
"os"
"path/filepath"
)

// Config holds all application configuration
type Config struct {
App      AppConfig
Database DatabaseConfig
AI      AIConfig
Server  ServerConfig
}

// AppConfig holds application settings
type AppConfig struct {
Name        string
Version     string
Environment string
}

// DatabaseConfig holds database settings
type DatabaseConfig struct {
Path string
}

// AIConfig holds AI provider settings
type AIConfig struct {
Provider     string
OllamaHost   string
OllamaModel  string
GeminiAPIKey string
GeminiModel  string
}

// ServerConfig holds server settings
type ServerConfig struct {
Host string
Port int
}

func getEnv(key, defaultVal string) string {
if v := os.Getenv(key); v != "" {
return v
}
return defaultVal
}

// Load loads configuration from environment
func Load() *Config {
home, _ := os.UserHomeDir()
configDir := filepath.Join(home, ".ai-desktop-assistant")
os.MkdirAll(configDir, 0755)

return &Config{
App: AppConfig{
Name:        "AI Desktop Assistant",
Version:     "0.1.0",
Environment: getEnv("ENVIRONMENT", "dev"),
},
Database: DatabaseConfig{
Path: getEnv("DATABASE_PATH", filepath.Join(configDir, "app.db")),
},
AI: AIConfig{
Provider:    getEnv("AI_PROVIDER", "gemini"),
OllamaHost:  getEnv("OLLAMA_HOST", "http://localhost:11434"),
OllamaModel: getEnv("OLLAMA_MODEL", "qwen2.5-coder"),
GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
GeminiModel:  getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
},
Server: ServerConfig{
Host: getEnv("SERVER_HOST", "localhost"),
Port: 8080,
},
}
}
