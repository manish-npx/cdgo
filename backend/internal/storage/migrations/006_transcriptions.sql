CREATE TABLE IF NOT EXISTS transcriptions (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    text TEXT NOT NULL,
    language TEXT NOT NULL DEFAULT 'en',
    duration INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
