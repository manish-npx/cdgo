package main

import (
"context"
"embed"
"log"
"net/http"
"os"
"os/signal"
"syscall"

"ai-desktop-assistant/backend/internal/config"
"ai-desktop-assistant/backend/internal/handlers"
"ai-desktop-assistant/backend/internal/services/ai"
"ai-desktop-assistant/backend/internal/services/storage"

"github.com/wailsapp/wails/v2"
"github.com/wailsapp/wails/v2/pkg/options"
"github.com/wailsapp/wails/v2/pkg/options/assetserver"
"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
wails.App
config        *config.Config
storageService *storage.StorageService
aiService     *ai.AIService
handlers      *handlers.Handlers
}

func NewApp() *App {
app := &App{}
app.config = config.Load()
app.storageService = storage.New(app.config.StoragePath)
app.aiService = ai.NewGeminiService()
app.handlers = handlers.New(app.storageService, app.aiService, app.config)
return app
}

func (a *App) startup(ctx context.Context) {
log.Println("AI Desktop Assistant starting...")
if err := a.storageService.Init(); err != nil {
log.Printf("Storage error: %v", err)
}
log.Println("Application ready")
}

func (a *App) shutdown(ctx context.Context) {
log.Println("Shutting down...")
a.storageService.Close()
}

// =====================
// Wails Bindings - API Key Management
// =====================

func (a *App) GetAPIKey(ctx context.Context) string {
return a.config.GeminiAPIKey
}

func (a *App) SetAPIKey(ctx context.Context, key string) error {
a.config.GeminiAPIKey = key
a.aiService.SetAPIKey(key)
return a.config.Save()
}

func (a *App) GetAIModel(ctx context.Context) string {
return a.config.AIModel
}

func (a *App) SetAIModel(ctx context.Context, model string) error {
a.config.AIModel = model
return a.config.Save()
}

// =====================
// Wails Bindings - AI Chat
// =====================

func (a *App) SendMessage(ctx context.Context, message string) (string, error) {
if a.config.GeminiAPIKey == "" {
return "", nil
}
return a.aiService.Chat(ctx, message)
}

func (a *App) GetChatHistory(ctx context.Context) []map[string]interface{} {
return a.storageService.GetSessions()
}

// =====================
// Wails Bindings - Settings
// =====================

func (a *App) GetSettings(ctx context.Context) map[string]interface{} {
return map[string]interface{}{
"apiKey":       a.config.GeminiAPIKey,
"aiModel":      a.config.AIModel,
"overlayWidth": a.config.OverlayWidth,
"overlayHeight":a.config.OverlayHeight,
"alwaysOnTop":  a.config.AlwaysOnTop,
"darkMode":     a.config.DarkMode,
}
}

func (a *App) SaveSettings(ctx context.Context, settings map[string]interface{}) error {
if key, ok := settings["apiKey"].(string); ok {
a.config.GeminiAPIKey = key
a.aiService.SetAPIKey(key)
}
if model, ok := settings["aiModel"].(string); ok {
a.config.AIModel = model
}
if w, ok := settings["overlayWidth"].(float64); ok {
a.config.OverlayWidth = int(w)
}
if h, ok := settings["overlayHeight"].(float64); ok {
a.config.OverlayHeight = int(h)
}
if top, ok := settings["alwaysOnTop"].(bool); ok {
a.config.AlwaysOnTop = top
}
if dark, ok := settings["darkMode"].(bool); ok {
a.config.DarkMode = dark
}
return a.config.Save()
}

// =====================
// Wails Bindings - Overlay
// =====================

func (a *App) SetOverlaySize(ctx context.Context, width, height int) {
a.config.OverlayWidth = width
a.config.OverlayHeight = height
a.config.Save()
}

func (a *App) SetAlwaysOnTop(ctx context.Context, enabled bool) {
a.config.AlwaysOnTop = enabled
a.config.Save()
}

func main() {
app := NewApp()

wailsConfig := &options.App{
Title:  "AI Desktop Assistant",
Width:  app.config.OverlayWidth,
Height: app.config.OverlayHeight,
AssetServer: &assetserver.Options{
Assets: assets,
},
BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
OnStartup:  app.startup,
OnShutdown: app.shutdown,
Bind: []interface{}{app},
}

// Platform-specific options
wailsConfig.Windows = &windows.Options{
AlwaysOnTop: app.config.AlwaysOnTop,
}

err := wails.Run(wailsConfig)
if err != nil {
log.Fatal("Error:", err)
}
}
