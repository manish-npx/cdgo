package main

import (
"context"
"embed"
"log"

"ai-desktop-assistant/backend/internal/config"
"ai-desktop-assistant/backend/internal/services/ai"
"ai-desktop-assistant/backend/internal/storage"

"github.com/wailsapp/wails/v2"
"github.com/wailsapp/wails/v2/pkg/options"
"github.com/wailsapp/wails/v2/pkg/options/assetserver"
"github.com/wailsapp/wails/v2/pkg/options/windows"
"go.uber.org/zap"
)

//go:embed all:frontend/dist
var assets embed.FS

// App is the main application struct
type App struct {
wails.App
config     *config.Config
store      *storage.Store
aiService  *ai.AIService
logger     *zap.Logger
}

// NewApp creates a new application instance
func NewApp() *App {
logger, _ := zap.NewDevelopment()
defer logger.Sync()

cfg := config.Load()

// Initialize database
db, err := storage.New(cfg.Database.Path, logger)
if err != nil {
logger.Fatal("Failed to initialize database", zap.Error(err))
}

store := storage.NewStore(db)
aiService := ai.NewService(cfg)

logger.Info("Application initialized", zap.String("version", cfg.App.Version))

return &App{
config:    cfg,
store:     store,
aiService: aiService,
logger:   logger,
}
}

// startup is called when the application starts
func (a *App) startup(ctx context.Context) {
a.logger.Info("Application starting...")
}

// shutdown is called when the application closes
func (a *App) shutdown(ctx context.Context) {
a.logger.Info("Application shutting down...")
}

// ====================
// Wails Bindings - Settings & API Key
// ====================

// GetSettings returns current application settings
func (a *App) GetSettings(ctx context.Context) map[string]interface{} {
// Load from database
apiKey, _ := a.store.Config.Get("gemini_api_key")
geminiModel, _ := a.store.Config.Get("gemini_model")
ollamaHost, _ := a.store.Config.Get("ollama_host")
ollamaModel, _ := a.store.Config.Get("ollama_model")
aiProvider, _ := a.store.Config.Get("ai_provider")

if apiKey == "" {
apiKey = a.config.AI.GeminiAPIKey
}
if geminiModel == "" {
geminiModel = a.config.AI.GeminiModel
}
if ollamaHost == "" {
ollamaHost = a.config.AI.OllamaHost
}
if ollamaModel == "" {
ollamaModel = a.config.AI.OllamaModel
}
if aiProvider == "" {
aiProvider = a.config.AI.Provider
}

return map[string]interface{}{
"apiKey":      apiKey,
"geminiModel": geminiModel,
"ollamaHost":  ollamaHost,
"ollamaModel": ollamaModel,
"aiProvider":  aiProvider,
}
}

// SaveSettings saves application settings
func (a *App) SaveSettings(ctx context.Context, settings map[string]interface{}) error {
if apiKey, ok := settings["apiKey"].(string); ok {
a.store.Config.Set("gemini_api_key", apiKey)
a.aiService.SetAPIKey(apiKey)
}
if model, ok := settings["geminiModel"].(string); ok {
a.store.Config.Set("gemini_model", model)
a.aiService.SetModel(model)
}
if provider, ok := settings["aiProvider"].(string); ok {
a.store.Config.Set("ai_provider", provider)
a.aiService.SetProvider(provider)
}
if host, ok := settings["ollamaHost"].(string); ok {
a.store.Config.Set("ollama_host", host)
}
if model, ok := settings["ollamaModel"].(string); ok {
a.store.Config.Set("ollama_model", model)
}
return nil
}

// ====================
// Wails Bindings - AI Chat
// ====================

// SendMessage sends a message to the AI and returns the response
func (a *App) SendMessage(ctx context.Context, message string) (string, error) {
apiKey, _ := a.store.Config.Get("gemini_api_key")
if apiKey == "" {
return "", nil
}
a.aiService.SetAPIKey(apiKey)
return a.aiService.Chat(ctx, message)
}

// GetChatHistory returns chat history for a session
func (a *App) GetChatHistory(ctx context.Context) []map[string]interface{} {
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

// CreateSession creates a new chat session
func (a *App) CreateSession(ctx context.Context) (string, error) {
session, err := a.store.Session.Create("New Chat")
if err != nil {
return "", err
}
return session.ID, nil
}

// GetMessages returns messages for a session
func (a *App) GetMessages(ctx context.Context, sessionID string) []map[string]interface{} {
messages, _ := a.store.Message.GetBySession(sessionID)
result := make([]map[string]interface{}, len(messages))
for i, m := range messages {
result[i] = map[string]interface{}{
"id":    m.ID,
"role":  m.Role,
"text":  m.Content,
"date":  m.CreatedAt.Format("15:04"),
}
}
return result
}

// SaveMessage saves a message to the database
func (a *App) SaveMessage(ctx context.Context, sessionID, role, content string) error {
_, err := a.store.Message.Create(sessionID, role, content)
return err
}

// ====================
// Main Entry Point
// ====================

func main() {
app := NewApp()

wailsConfig := &options.App{
Title:  "AI Desktop Assistant",
Width:  420,
Height: 600,
AssetServer: &assetserver.Options{
Assets: assets,
},
BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
OnStartup:  app.startup,
OnShutdown: app.shutdown,
Bind: []interface{}{app},
Windows: &windows.Options{
AlwaysOnTop: true,
},
}

err := wails.Run(wailsConfig)
if err != nil {
log.Fatal("Error running application:", err)
}
}
