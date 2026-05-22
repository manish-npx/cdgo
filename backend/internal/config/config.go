package config

import (
"os"
"path/filepath"
"sync"

"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
App      AppConfig      `json:"app"`
Database DatabaseConfig `json:"database"`
AI       AIConfig       `json:"ai"`
Server   ServerConfig   `json:"server"`
mu       sync.RWMutex
}

type AppConfig struct {
Name        string `json:"name"`
Version     string `json:"version"`
Environment string `json:"environment"`
}

type DatabaseConfig struct {
Path         string `json:"path"`
MaxOpenConns int    `json:"maxOpenConns"`
MaxIdleConns int    `json:"maxIdleConns"`
}

type AIConfig struct {
Provider     string `json:"provider"`
OllamaHost   string `json:"ollamaHost"`
OllamaModel  string `json:"ollamaModel"`
GeminiAPIKey string `json:"geminiApiKey"`
GeminiModel  string `json:"geminiModel"`
OpenAIKey    string `json:"openAiKey"`
OpenAIModel  string `json:"openAiModel"`
GroqAPIKey   string `json:"groqApiKey"`
GroqModel    string `json:"groqModel"`
}

type ServerConfig struct {
Host string `json:"host"`
Port int    `json:"port"`
}

var globalConfig *Config

// Load loads configuration from environment variables and .env file
func Load() *Config {
godotenv.Load()

cfg := &Config{
App: AppConfig{
Name:        getEnv("APP_NAME", "AI Desktop Assistant"),
Version:     getEnv("APP_VERSION", "0.1.0"),
Environment: getEnv("ENVIRONMENT", "dev"),
},
Database: DatabaseConfig{
Path:         getEnv("DATABASE_PATH", "./data/app.db"),
MaxOpenConns: 25,
MaxIdleConns: 5,
},
AI: AIConfig{
Provider:    getEnv("AI_PROVIDER", "gemini"),
OllamaHost:  getEnv("OLLAMA_HOST", "http://localhost:11434"),
OllamaModel: getEnv("OLLAMA_MODEL", "qwen2.5-coder"),
GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
GeminiModel:  getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
OpenAIKey:    getEnv("OPENAI_API_KEY", ""),
OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-4"),
GroqAPIKey:    getEnv("GROQ_API_KEY", ""),
GroqModel:     getEnv("GROQ_MODEL", "llama-3.1-70b-versatile"),
},
Server: ServerConfig{
Host: getEnv("SERVER_HOST", "localhost"),
Port: 8080,
},
}

// Ensure database directory exists
if cfg.Database.Path != "" {
dir := filepath.Dir(cfg.Database.Path)
os.MkdirAll(dir, 0755)
}

globalConfig = cfg
return cfg
}

func getEnv(key, defaultValue string) string {
if value := os.Getenv(key); value != "" {
return value
}
return defaultValue
}

func Get() *Config {
if globalConfig == nil {
return Load()
}
return globalConfig
}

// Update updates a specific configuration value
func (c *Config) Update(key string, value string) error {
c.mu.Lock()
defer c.mu.Unlock()

switch key {
case "ai.geminiApiKey":
c.AI.GeminiAPIKey = value
case "ai.geminiModel":
c.AI.GeminiModel = value
case "ai.ollamaHost":
c.AI.OllamaHost = value
case "ai.ollamaModel":
c.AI.OllamaModel = value
default:
// Store in config file
}

// Save to database
return nil
}
