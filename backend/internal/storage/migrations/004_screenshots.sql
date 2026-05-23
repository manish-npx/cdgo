CREATE TABLE IF NOT EXISTS screenshots (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    path TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
