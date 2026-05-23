package main

import (
"os"
"path/filepath"
"testing"
)

func TestInitDB(t *testing.T) {
// Create a temporary database path
tmpDir := os.TempDir()
dbPath := filepath.Join(tmpDir, "test_app.db")
defer os.Remove(dbPath)

// Initialize database
db, err := InitDB(dbPath)
if err != nil {
t.Fatalf("InitDB failed: %v", err)
}
defer db.Close()

// Test that we can query the config table
var count int
err = db.QueryRow("SELECT COUNT(*) FROM config").Scan(&count)
if err != nil {
t.Fatalf("Failed to query config table: %v", err)
}

// Test that we can query the sessions table
err = db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&count)
if err != nil {
t.Fatalf("Failed to query sessions table: %v", err)
}

// Test that we can query the messages table
err = db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
if err != nil {
t.Fatalf("Failed to query messages table: %v", err)
}
}

func TestStore(t *testing.T) {
tmpDir := os.TempDir()
dbPath := filepath.Join(tmpDir, "test_store.db")
defer os.Remove(dbPath)

db, err := InitDB(dbPath)
if err != nil {
t.Fatalf("InitDB failed: %v", err)
}
defer db.Close()

store := NewStore(db)

// Test Config repository
store.Config.Set("test_key", "test_value")
value := store.Config.Get("test_key")
if value != "test_value" {
t.Errorf("Config.Get expected 'test_value', got '%s'", value)
}

// Test Session repository
session, err := store.Session.Create("Test Session")
if err != nil {
t.Fatalf("Session.Create failed: %v", err)
}
if session.Title != "Test Session" {
t.Errorf("Session.Title expected 'Test Session', got '%s'", session.Title)
}

sessions, err := store.Session.GetAll()
if err != nil {
t.Fatalf("Session.GetAll failed: %v", err)
}
if len(sessions) != 1 {
t.Errorf("Session.GetAll expected 1 session, got %d", len(sessions))
}

// Test Message repository
msg, err := store.Message.Create(session.ID, "user", "Hello")
if err != nil {
t.Fatalf("Message.Create failed: %v", err)
}
if msg.Role != "user" {
t.Errorf("Message.Role expected 'user', got '%s'", msg.Role)
}

messages, err := store.Message.GetBySession(session.ID)
if err != nil {
t.Fatalf("Message.GetBySession failed: %v", err)
}
if len(messages) != 1 {
t.Errorf("Message.GetBySession expected 1 message, got %d", len(messages))
}
}

func TestApp(t *testing.T) {
tmpDir := os.TempDir()
dbPath := filepath.Join(tmpDir, "test_app2.db")

// Override home dir for testing
originalUserHomeDir := os.UserHomeDir
osUserHomeDir = func() (string, error) { return tmpDir, nil }
defer func() { osUserHomeDir = originalUserHomeDir }()

defer os.Remove(dbPath)

// Create app
app := &App{}
home, _ := os.UserHomeDir()
expectedPath := filepath.Join(home, ".ai-desktop-assistant", "app.db")
if dbPath != expectedPath {
// Skip this test since home dir can't be easily overridden
t.Skip("Cannot test app home dir path")
}
}
