# AI Desktop Assistant

A production-grade Windows desktop AI assistant built with **Go + Wails + React** for educational purposes.

## Features

- **Native Windows Overlay** - Always-on-top desktop window
- **Multiple AI Providers** - Gemini API, Ollama (local), OpenAI, Groq
- **SQLite Database** - Sessions and messages stored locally
- **Real-time Streaming** - WebSocket-based AI responses
- **Event Bus Architecture** - Decoupled component communication
- **OCR Pipeline** - Screenshot capture and text extraction
- **Audio Transcription** - Voice-to-text with whisper.cpp

## Tech Stack

| Layer | Technology |
|-------|------------|
| Desktop | Wails v2 |
| Backend | Go 1.24+ |
| Frontend | React 19 + TypeScript |
| Database | SQLite with migrations |
| AI | Ollama, Gemini, OpenAI, Groq |
| Logging | Uber Zap |

## Project Structure

```
cdgo/
├── backend/
│   ├── cmd/                  # Application entry
│   └── internal/
│       ├── ai/              # AI provider integration
│       ├── audio/           # Audio processing
│       ├── capture/          # Screenshot capture
│       ├── config/           # Configuration management
│       ├── events/           # Event bus system
│       ├── logging/          # Structured logging
│       ├── overlay/           # Overlay window controls
│       ├── storage/           # SQLite + repositories
│       └── websocket/         # Real-time messaging
├── frontend/                 # React UI
├── docs/                     # Documentation
└── scripts/                  # Build scripts
```

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Setup

```bash
# Clone repository
git clone https://github.com/manish-npx/cdgo.git
cd cdgo

# Install frontend dependencies
cd frontend && npm install && npm run build && cd ..

# Run development
wails dev
```

### Configuration

Create `.env` file:
```
GEMINI_API_KEY=your-api-key-here
AI_PROVIDER=gemini
```

## Documentation

See `docs/` folder for detailed documentation:
- `docs/architecture.md` - System design
- `docs/tdd-workflow.md` - Development workflow
- `docs/ai-usage.md` - AI provider setup

## License

MIT - Educational Use
