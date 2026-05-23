package storage

import (
"database/sql"
"embed"
"fmt"
"os"
"path/filepath"
"sort"
"sync"

_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Logger interface for storage
type Logger interface {
Info(string, ...interface{})
Error(string, ...interface{})
}

// Database is the SQLite database wrapper
type Database struct {
db     *sql.DB
logger Logger
mu     sync.RWMutex
}

// New creates a new database connection
func New(dbPath string, logger Logger) (*Database, error) {
dir := filepath.Dir(dbPath)
if err := os.MkdirAll(dir, 0755); err != nil {
return nil, fmt.Errorf("failed to create directory: %w", err)
}

db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_busy_timeout=5000")
if err != nil {
return nil, fmt.Errorf("failed to open database: %w", err)
}

db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)

database := &Database{db: db, logger: logger}

if err := database.Migrate(); err != nil {
db.Close()
return nil, fmt.Errorf("failed to migrate: %w", err)
}

if logger != nil {
logger.Info("Database initialized", "path", dbPath)
}

return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
d.mu.Lock()
defer d.mu.Unlock()
return d.db.Close()
}

// GetDB returns the underlying sql.DB
func (d *Database) GetDB() *sql.DB {
d.mu.RLock()
defer d.mu.RUnlock()
return d.db
}

// WithTx executes a function within a transaction
func (d *Database) WithTx(fn func(*sql.Tx) error) error {
tx, err := d.db.Begin()
if err != nil {
return fmt.Errorf("failed to begin transaction: %w", err)
}

if err := fn(tx); err != nil {
tx.Rollback()
return err
}

return tx.Commit()
}

// Migrate runs pending database migrations
func (d *Database) Migrate() error {
entries, err := migrationsFS.ReadDir("migrations")
if err != nil {
return fmt.Errorf("failed to read migrations: %w", err)
}

var files []string
for _, entry := range entries {
if !entry.IsDir() {
files = append(files, entry.Name())
}
}
sort.Strings(files)

// Create migrations table
_, _ = d.db.Exec(`
CREATE TABLE IF NOT EXISTS migrations (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL UNIQUE,
applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
`)

for _, filename := range files {
var count int
_ = d.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", filename).Scan(&count)
if count > 0 {
continue
}

content, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
if err != nil {
return fmt.Errorf("failed to read %s: %w", filename, err)
}

_, err = d.db.Exec(string(content))
if err != nil {
return fmt.Errorf("failed to execute %s: %w", filename, err)
}

_, _ = d.db.Exec("INSERT INTO migrations (name) VALUES (?)", filename)

if d.logger != nil {
d.logger.Info("Applied migration", "file", filename)
}
}

return nil
}

// Health checks database connectivity
func (d *Database) Health() error {
return d.db.Ping()
}
