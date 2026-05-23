package main

import (
"database/sql"
"fmt"
"os"
"path/filepath"
"time"

"github.com/google/uuid"
_ "github.com/mattn/go-sqlite3"
)

// ====================
// Config Repository
// ====================

type ConfigRepository struct {
db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
return &ConfigRepository{db: db}
}

func (r *ConfigRepository) Get(key string) string {
var value string
r.db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
return value
}

func (r *ConfigRepository) Set(key, value string) {
r.db.Exec(`
INSERT INTO config (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
`, key, value, value)
}

// ====================
// Session Repository
// ====================

type Session struct {
ID        string
Title     string
CreatedAt time.Time
UpdatedAt time.Time
}

type SessionRepository struct {
db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(title string) (*Session, error) {
session := &Session{
ID:        uuid.New().String(),
Title:     title,
CreatedAt: time.Now(),
UpdatedAt: time.Now(),
}
_, err := r.db.Exec(`
INSERT INTO sessions (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)
`, session.ID, session.Title, session.CreatedAt, session.UpdatedAt)
return session, err
}

func (r *SessionRepository) GetAll() ([]Session, error) {
rows, err := r.db.Query(`
SELECT id, title, created_at, updated_at FROM sessions ORDER BY updated_at DESC
`)
if err != nil {
return nil, err
}
defer rows.Close()

var sessions []Session
for rows.Next() {
var s Session
if err := rows.Scan(&s.ID, &s.Title, &s.CreatedAt, &s.UpdatedAt); err != nil {
return nil, err
}
sessions = append(sessions, s)
}
return sessions, rows.Err()
}

func (r *SessionRepository) Get(id string) (*Session, error) {
var s Session
err := r.db.QueryRow(`
SELECT id, title, created_at, updated_at FROM sessions WHERE id = ?
`, id).Scan(&s.ID, &s.Title, &s.CreatedAt, &s.UpdatedAt)
if err != nil {
return nil, err
}
return &s, nil
}

func (r *SessionRepository) Delete(id string) error {
_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", id)
return err
}

// ====================
// Message Repository
// ====================

type Message struct {
ID        string
SessionID string
Role      string
Content   string
CreatedAt time.Time
}

type MessageRepository struct {
db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(sessionID, role, content string) (*Message, error) {
msg := &Message{
ID:        uuid.New().String(),
SessionID: sessionID,
Role:      role,
Content:   content,
CreatedAt: time.Now(),
}
_, err := r.db.Exec(`
INSERT INTO messages (id, session_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)
`, msg.ID, msg.SessionID, msg.Role, msg.Content, msg.CreatedAt)
return msg, err
}

func (r *MessageRepository) GetBySession(sessionID string) ([]Message, error) {
rows, err := r.db.Query(`
SELECT id, session_id, role, content, created_at FROM messages
WHERE session_id = ? ORDER BY created_at ASC
`, sessionID)
if err != nil {
return nil, err
}
defer rows.Close()

var messages []Message
for rows.Next() {
var m Message
if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
return nil, err
}
messages = append(messages, m)
}
return messages, rows.Err()
}

// ====================
// Store Facade
// ====================

type Store struct {
Config  *ConfigRepository
Session *SessionRepository
Message *MessageRepository
}

func NewStore(db *sql.DB) *Store {
return &Store{
Config:  NewConfigRepository(db),
Session: NewSessionRepository(db),
Message: NewMessageRepository(db),
}
}

// ====================
// Database Init
// ====================

func InitDB(dbPath string) (*sql.DB, error) {
dir := filepath.Dir(dbPath)
os.MkdirAll(dir, 0755)

db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL")
if err != nil {
return nil, err
}

// Create tables
tables := []string{
`CREATE TABLE IF NOT EXISTS config (
key TEXT PRIMARY KEY,
value TEXT NOT NULL,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)`,
`CREATE TABLE IF NOT EXISTS sessions (
id TEXT PRIMARY KEY,
title TEXT NOT NULL,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)`,
`CREATE TABLE IF NOT EXISTS messages (
id TEXT PRIMARY KEY,
session_id TEXT NOT NULL,
role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
content TEXT NOT NULL,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
)`,
`CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id)`,
}

for _, table := range tables {
if _, err := db.Exec(table); err != nil {
return nil, err
}
}

return db, nil
}

// ====================
// App
// ====================

type App struct {
store *Store
}

func NewApp() *App {
home, _ := os.UserHomeDir()
dbPath := filepath.Join(home, ".ai-desktop-assistant", "app.db")

db, err := InitDB(dbPath)
if err != nil {
fmt.Println("❌ Database error:", err)
os.Exit(1)
}

return &App{store: NewStore(db)}
}

func (a *App) GetSettings() map[string]interface{} {
return map[string]interface{}{
"apiKey":      a.store.Config.Get("gemini_api_key"),
"geminiModel": a.store.Config.Get("gemini_model"),
"aiProvider":  a.store.Config.Get("ai_provider"),
}
}

func (a *App) SaveSettings(settings map[string]interface{}) string {
for key, val := range settings {
if str, ok := val.(string); ok {
a.store.Config.Set(key, str)
}
}
return "Settings saved!"
}

func (a *App) SendMessage(message string, sessionID string) string {
apiKey := a.store.Config.Get("gemini_api_key")
if apiKey == "" {
return "⚠️ Please configure your API key in Settings!"
}

a.store.Message.Create(sessionID, "user", message)
a.store.Message.Create(sessionID, "assistant", "AI response placeholder - integrate AI service here")

return "Response from AI (placeholder)"
}

func (a *App) CreateSession() string {
session, _ := a.store.Session.Create("New Chat")
return session.ID
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
// Main
// ====================

func main() {
app := NewApp()

fmt.Println(`
╔════════════════════════════════════════════╗
║   🤖 AI Desktop Assistant                  ║
║   Go + SQLite + Gemini                      ║
╚════════════════════════════════════════════╝
`)

fmt.Println("✅ Database initialized!")

settings := app.GetSettings()
fmt.Printf("\n⚙️  Settings: %v\n", settings)

if settings["apiKey"] == "" {
fmt.Println("\n⚠️  No API Key configured!")
fmt.Println("\nTo configure, create .env file:")
fmt.Println("   GEMINI_API_KEY=your-api-key-here")
fmt.Println("\nGet free key: https://aistudio.google.com/app/apikey")
} else {
fmt.Println("\n✅ API Key configured!")

// Demo
sessionID := app.CreateSession()
fmt.Println("\n💬 Demo: Creating session...", sessionID)

response := app.SendMessage("Hello!", sessionID)
fmt.Println("📝 AI:", response)
}

home, _ := os.UserHomeDir()
fmt.Println("\n📁 Data stored at:", filepath.Join(home, ".ai-desktop-assistant"))
}
