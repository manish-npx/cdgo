package storage

import (
"encoding/json"
"os"
"path/filepath"
"sync"
"time"
)

type StorageService struct {
storagePath string
mu          sync.RWMutex
sessions    []Session
}

type Session struct {
ID        string    `json:"id"`
Timestamp time.Time `json:"timestamp"`
UserMsg   string    `json:"userMsg"`
AIMsg     string    `json:"aiMsg"`
Model     string    `json:"model"`
}

func New(storagePath string) *StorageService {
dir := filepath.Dir(storagePath)
os.MkdirAll(dir, 0755)
return &StorageService{storagePath: storagePath}
}

func (s *StorageService) Init() error {
return s.load()
}

func (s *StorageService) load() error {
s.mu.Lock()
defer s.mu.Unlock()

data, err := os.ReadFile(s.storagePath)
if err != nil {
return nil
}

return json.Unmarshal(data, &s.sessions)
}

func (s *StorageService) Save() error {
s.mu.Lock()
defer s.mu.Unlock()

dir := filepath.Dir(s.storagePath)
os.MkdirAll(dir, 0755)

data, _ := json.MarshalIndent(s.sessions, "", "  ")
return os.WriteFile(s.storagePath, data, 0644)
}

func (s *StorageService) Close() error {
return s.Save()
}

func (s *StorageService) AddSession(session *Session) error {
s.mu.Lock()
defer s.mu.Unlock()

if session.ID == "" {
session.ID = time.Now().Format("20060102150405")
}
session.Timestamp = time.Now()

s.sessions = append([]Session{*session}, s.sessions...)
if len(s.sessions) > 100 {
s.sessions = s.sessions[:100]
}

return s.Save()
}

func (s *StorageService) GetSessions() []map[string]interface{} {
s.mu.RLock()
defer s.mu.RUnlock()

result := make([]map[string]interface{}, len(s.sessions))
for i, session := range s.sessions {
result[i] = map[string]interface{}{
"id":        session.ID,
"timestamp": session.Timestamp.Format(time.RFC3339),
"userMsg":   session.UserMsg,
"aiMsg":     session.AIMsg,
"model":     session.Model,
}
}

if result == nil {
return []map[string]interface{}{}
}
return result
}

func (s *StorageService) ClearSessions() error {
s.mu.Lock()
defer s.mu.Unlock()

s.sessions = []Session{}
return s.Save()
}
