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

type Database struct {
db     *sql.DB
logger Logger
mu     sync.RWMutex
}

type Logger interface {
Info(string, ...interface{})
Error(string, ...interface{})
Debug(string, ...interface{})
}

func New(dbPath string, logger Logger) (*Database, error) {
dir := filepath.Dir(dbPath)
if err := os.MkdirAll(dir, 0755); err != nil {
return nil, fmt.Errorf("failed to create database directory: %w", err)
}

db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_busy_timeout=5000")
if err != nil {
return nil, fmt.Errorf("failed to open database: %w", err)
}

db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)

database := &Database{
db:     db,
logger: logger,
}

if err := database.Migrate(); err != nil {
db.Close()
return nil, fmt.Errorf("failed to run migrations: %w", err)
}

logger.Info("Database initialized at %s", dbPath)
return database, nil
}

func (d *Database) Close() error {
d.mu.Lock()
defer d.mu.Unlock()
return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
d.mu.RLock()
defer d.mu.RUnlock()
return d.db
}

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

func (d *Database) Migrate() error {
entries, err := migrationsFS.ReadDir("migrations")
if err != nil {
return fmt.Errorf("failed to read migrations directory: %w", err)
}

var migrationFiles []string
for _, entry := range entries {
if !entry.IsDir() {
migrationFiles = append(migrationFiles, entry.Name())
}
}
sort.Strings(migrationFiles)

// Create migrations table if not exists
_, err = d.db.Exec(`
CREATE TABLE IF NOT EXISTS migrations (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL UNIQUE,
applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
`)
if err != nil {
return fmt.Errorf("failed to create migrations table: %w", err)
}

for _, filename := range migrationFiles {
// Check if migration already applied
var count int
err := d.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", filename).Scan(&count)
if err != nil {
return fmt.Errorf("failed to check migration: %w", err)
}

if count > 0 {
continue
}

// Read and execute migration
sqlContent, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
if err != nil {
return fmt.Errorf("failed to read migration %s: %w", filename, err)
}

_, err = d.db.Exec(string(sqlContent))
if err != nil {
return fmt.Errorf("failed to execute migration %s: %w", filename, err)
}

// Record migration
_, err = d.db.Exec("INSERT INTO migrations (name) VALUES (?)", filename)
if err != nil {
return fmt.Errorf("failed to record migration %s: %w", filename, err)
}

d.logger.Info("Applied migration: %s", filename)
}

return nil
}

func (d *Database) Health() error {
return d.db.Ping()
}
