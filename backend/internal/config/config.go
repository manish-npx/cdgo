package config

import (
"encoding/json"
"os"
"path/filepath"
)

type Config struct {
StoragePath   string `json:"storagePath"`
GeminiAPIKey  string `json:"geminiApiKey"`
AIModel       string `json:"aiModel"`
OverlayWidth  int    `json:"overlayWidth"`
OverlayHeight int    `json:"overlayHeight"`
AlwaysOnTop   bool   `json:"alwaysOnTop"`
DarkMode      bool   `json:"darkMode"`
}

func DefaultConfig() *Config {
homeDir, _ := os.UserHomeDir()
configDir := filepath.Join(homeDir, ".ai-desktop-assistant")
os.MkdirAll(configDir, 0755)

return &Config{
StoragePath:   filepath.Join(configDir, "data.db"),
GeminiAPIKey:  "",
AIModel:       "gemini-2.0-flash",
OverlayWidth:  400,
OverlayHeight: 500,
AlwaysOnTop:   true,
DarkMode:      true,
}
}

func Load() *Config {
cfg := DefaultConfig()
homeDir, _ := os.UserHomeDir()
configDir := filepath.Join(homeDir, ".ai-desktop-assistant")
configPath := filepath.Join(configDir, "config.json")

data, err := os.ReadFile(configPath)
if err != nil {
return cfg
}

json.Unmarshal(data, cfg)
return cfg
}

func (c *Config) Save() error {
homeDir, _ := os.UserHomeDir()
configDir := filepath.Join(homeDir, ".ai-desktop-assistant")
configPath := filepath.Join(configDir, "config.json")
os.MkdirAll(configDir, 0755)

data, _ := json.MarshalIndent(c, "", "  ")
return os.WriteFile(configPath, data, 0644)
}
