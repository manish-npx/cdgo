package main

import (
"fmt"
"log"
"os"

"ai-desktop-assistant/backend/internal/config"
"ai-desktop-assistant/backend/internal/logging"
"ai-desktop-assistant/backend/internal/storage"
)

// App is the main application struct
type App struct {
config  *config.Config
logger  *logging.Logger
store   *storage.Store
}

// NewApp creates a new application instance
func NewApp() *App {
logger := logging.New()
logger.Info("Starting AI Desktop Assistant...")

cfg := config.Load()
logger.Info("Configuration loaded", "version", cfg.App.Version)

store, err := storage.New(cfg.Database.Path, logger)
if err != nil {
logger.Error("Failed to initialize storage", "error", err)
os.Exit(1)
}
logger.Info("Storage initialized", "path", cfg.Database.Path)

return &App{
config: cfg,
logger: logger,
store:  store,
}
}

// ====================
// Methods for Settings
// ====================

func (a *App) GetSettings() map[string]interface{} {
return map[string]interface{}{
"apiKey":      a.store.Config.Get("gemini_api_key"),
"geminiModel": a.store.Config.Get("gemini_model"),
"ollamaHost":  a.store.Config.Get("ollama_host"),
"ollamaModel": a.store.Config.Get("ollama_model"),
"aiProvider":  a.store.Config.Get("ai_provider"),
}
}

func (a *App) SaveSettings(settings map[string]interface{}) string {
for key, val := range settings {
if str, ok := val.(string); ok {
a.store.Config.Set(key, str)
}
}
return "Settings saved"
}

// ====================
// Methods for AI Chat
// ====================

func (a *App) SendMessage(message string, sessionID string) string {
apiKey := a.store.Config.Get("gemini_api_key")
if apiKey == "" {
return "Please configure your API key in Settings"
}

a.store.Message.Create(sessionID, "user", message)
response := "AI response placeholder - integrate AI service here"
a.store.Message.Create(sessionID, "assistant", response)

return response
}

func (a *App) GetChatHistory() []map[string]interface{} {
sessions, _ := a.store.Session.GetAll()
result := make([]map[string]interface{}, len(sessions))
for i, s := range sessions {
result[i] = map[string]interface{}{
"id":    s.ID,
"title": s.Title,
"date":  s.CreatedAt.Format("2006-01-02"),
}
}
return result
}

func (a *App) CreateSession() string {
session, _ := a.store.Session.Create("New Chat")
return session.ID
}

func (a *App) GetMessages(sessionID string) []map[string]interface{} {
messages, _ := a.store.Message.GetBySession(sessionID)
result := make([]map[string]interface{}, len(messages))
for i, m := range messages {
result[i] = map[string]interface{}{
"id":   m.ID,
"role": m.Role,
"text": m.Content,
"date": m.CreatedAt.Format("15:04"),
}
}
return result
}

// ====================
// Main Entry Point
// ====================

func main() {
app := NewApp()

fmt.Println(`
╔════════════════════════════════════════════╗
║   🤖 AI Desktop Assistant                  ║
║   Go + Wails Desktop Application           ║
╚════════════════════════════════════════════╝
`)

fmt.Printf("📁 Database: %s\n", app.config.Database.Path)
fmt.Printf("🤖 Model: %s\n\n", app.config.AI.GeminiModel)

settings := app.GetSettings()
fmt.Printf("⚙️  Settings: %v\n\n", settings)

if settings["apiKey"] == "" {
fmt.Println("⚠️  No API Key configured!")
fmt.Println("\nTo configure, create .env file with:")
fmt.Println("   GEMINI_API_KEY=your-api-key-here")
fmt.Println("\nGet free API key: https://aistudio.google.com/app/apikey")
}

fmt.Println(`
╔════════════════════════════════════════════╗
║   Ready for Wails Desktop App!            ║
║                                            ║
║   To build desktop app, install Wails:     ║
║   go install github.com/wailsapp/wails/v2/cmd/wails@latest
║   Then run: wails dev                      ║
╚════════════════════════════════════════════╝
`)
}
