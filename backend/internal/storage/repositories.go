package storage

import (
"database/sql"
"time"

"github.com/google/uuid"
)

// ConfigRepository handles key-value configuration
type ConfigRepository struct {
db *Database
}

func NewConfigRepository(db *Database) *ConfigRepository {
return &ConfigRepository{db: db}
}

func (r *ConfigRepository) Get(key string) string {
var value string
r.db.GetDB().QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
return value
}

func (r *ConfigRepository) Set(key, value string) {
r.db.GetDB().Exec(`
INSERT INTO config (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
`, key, value, value)
}

// SessionRepository handles chat sessions
type SessionRepository struct {
db *Database
}

type Session struct {
ID        string
Title     string
CreatedAt time.Time
UpdatedAt time.Time
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
var s Session
err := r.db.GetDB().QueryRow(`
SELECT id, title, created_at, updated_at FROM sessions WHERE id = ?
`, id).Scan(&s.ID, &s.Title, &s.CreatedAt, &s.UpdatedAt)

if err != nil {
return nil, err
}
return &s, nil
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

func (r *SessionRepository) Delete(id string) error {
_, err := r.db.GetDB().Exec("DELETE FROM sessions WHERE id = ?", id)
return err
}

// MessageRepository handles chat messages
type MessageRepository struct {
db *Database
}

type Message struct {
ID        string
SessionID string
Role      string
Content   string
CreatedAt time.Time
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
var m Message
err := r.db.GetDB().QueryRow(`
SELECT id, session_id, role, content, created_at FROM messages
WHERE session_id = ? ORDER BY created_at DESC LIMIT 1
`, sessionID).Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt)

if err == sql.ErrNoRows {
return nil, nil
}
return &m, err
}

// Store is the facade for all repositories
type Store struct {
Config  *ConfigRepository
Session *SessionRepository
Message *MessageRepository
}

func NewStore(db *Database) *Store {
return &Store{
Config:  NewConfigRepository(db),
Session: NewSessionRepository(db),
Message: NewMessageRepository(db),
}
}
