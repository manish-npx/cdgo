CREATE TABLE IF NOT EXISTS ocr_results (
    id TEXT PRIMARY KEY,
    screenshot_id TEXT,
    text TEXT NOT NULL,
    confidence REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (screenshot_id) REFERENCES screenshots(id) ON DELETE CASCADE
);
