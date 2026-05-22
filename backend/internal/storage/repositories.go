package storage

import (
"database/sql"
"fmt"
"time"

"github.com/google/uuid"
)

// ConfigRepository handles configuration data
type ConfigRepository struct {
db *Database
}

func NewConfigRepository(db *Database) *ConfigRepository {
return &ConfigRepository{db: db}
}

func (r *ConfigRepository) Get(key string) (string, error) {
var value string
err := r.db.GetDB().QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
if err == sql.ErrNoRows {
return "", nil
}
return value, err
}

func (r *ConfigRepository) Set(key, value string) error {
_, err := r.db.GetDB().Exec(`
INSERT INTO config (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
`, key, value, value)
return err
}

func (r *ConfigRepository) GetAll() (map[string]string, error) {
rows, err := r.db.GetDB().Query("SELECT key, value FROM config")
if err != nil {
return nil, err
}
defer rows.Close()

result := make(map[string]string)
for rows.Next() {
var key, value string
if err := rows.Scan(&key, &value); err != nil {
return nil, err
}
result[key] = value
}
return result, rows.Err()
}

// SessionRepository handles chat sessions
type SessionRepository struct {
db *Database
}

type Session struct {
ID        string    `json:"id"`
Title     string    `json:"title"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

func NewSessionRepository(db *Database) *SessionRepository {
return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(title string) (*Session, error) {
session := &Session{
ID:        uuid.New().String(),
Title:     title,
CreatedAt: time.Now(),
UpdatedAt: time.Now(),
}

_, err := r.db.GetDB().Exec(`
INSERT INTO sessions (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)
`, session.ID, session.Title, session.CreatedAt, session.UpdatedAt)

return session, err
}

func (r *SessionRepository) Get(id string) (*Session, error) {
var session Session
err := r.db.GetDB().QueryRow(`
SELECT id, title, created_at, updated_at FROM sessions WHERE id = ?
`, id).Scan(&session.ID, &session.Title, &session.CreatedAt, &session.UpdatedAt)

if err != nil {
return nil, err
}
return &session, nil
}

func (r *SessionRepository) GetAll() ([]Session, error) {
rows, err := r.db.GetDB().Query(`
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

func (r *SessionRepository) UpdateTitle(id, title string) error {
_, err := r.db.GetDB().Exec(`
UPDATE sessions SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
`, title, id)
return err
}

func (r *SessionRepository) Delete(id string) error {
_, err := r.db.GetDB().Exec("DELETE FROM sessions WHERE id = ?", id)
return err
}

// MessageRepository handles chat messages
type MessageRepository struct {
db *Database
}

type Message struct {
ID        string    `json:"id"`
SessionID string    `json:"session_id"`
Role      string    `json:"role"`
Content   string    `json:"content"`
CreatedAt time.Time `json:"created_at"`
}

func NewMessageRepository(db *Database) *MessageRepository {
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

_, err := r.db.GetDB().Exec(`
INSERT INTO messages (id, session_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)
`, msg.ID, msg.SessionID, msg.Role, msg.Content, msg.CreatedAt)

// Update session updated_at
r.db.GetDB().Exec("UPDATE sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = ?", sessionID)

return msg, err
}

func (r *MessageRepository) GetBySession(sessionID string) ([]Message, error) {
rows, err := r.db.GetDB().Query(`
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

func (r *MessageRepository) GetLast(sessionID string) (*Message, error) {
var msg Message
err := r.db.GetDB().QueryRow(`
SELECT id, session_id, role, content, created_at FROM messages 
WHERE session_id = ? ORDER BY created_at DESC LIMIT 1
`, sessionID).Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt)

if err != nil {
return nil, err
}
return &msg, nil
}

func (r *MessageRepository) Delete(id string) error {
_, err := r.db.GetDB().Exec("DELETE FROM messages WHERE id = ?", id)
return err
}

// Store facade for all repositories
type Store struct {
Config   *ConfigRepository
Session  *SessionRepository
Message  *MessageRepository
}

func NewStore(db *Database) *Store {
return &Store{
Config:  NewConfigRepository(db),
Session: NewSessionRepository(db),
Message: NewMessageRepository(db),
}
}

func (s *Store) Health() error {
return s.Config.db.Health()
}
