# AI Desktop Assistant 🤖

A production-grade desktop AI assistant built with **Go + Wails + React** for educational purposes.

## Features

- **Standalone Desktop App** - Native Windows/Mac/Linux application
- **Multiple AI Providers** - Gemini API, Ollama (local), OpenAI, Groq
- **SQLite Database** - Sessions and messages stored locally
- **Settings Panel** - Enter your API key and configure providers
- **Dark Theme UI** - Modern purple gradient interface
- **Session History** - Chat history persisted in SQLite

## 🚀 Quick Start

### Prerequisites

1. **Go 1.21+** - [Download](https://go.dev/dl/)
2. **Node.js 18+** - [Download](https://nodejs.org/)
3. **Wails CLI** - `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Setup

```bash
# Clone and enter project
git clone https://github.com/manish-npx/cdgo.git
cd cdgo/ai-desktop-assistant

# Install frontend
cd frontend
npm install
npm run build

# Run with Wails
cd ..
wails dev
```

### First Run

1. Click the **⚙️ Settings** button
2. Enter your **API key** (get free from [aistudio.google.com](https://aistudio.google.com/app/apikey))
3. Select **AI Provider** (Gemini recommended)
4. Click **Save**
5. Start chatting!

## 📁 Project Structure

```
ai-desktop-assistant/
├── backend/                    # Go + Wails backend
│   ├── cmd/main.go            # Application entry
│   └── internal/
│       ├── config/            # Configuration management
│       ├── services/ai/       # AI provider integration
│       └── storage/          # SQLite + repositories
│           └── migrations/   # Database migrations
├── frontend/                   # React UI
│   └── src/
│       └── App.tsx           # Main application
└── README.md
```

## 🔧 Configuration

Settings stored at: `~/.ai-desktop-assistant/config.json`

### Environment Variables

```bash
# AI Provider
AI_PROVIDER=gemini              # gemini, ollama, openai, groq
GEMINI_API_KEY=your-key
GEMINI_MODEL=gemini-2.0-flash
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=qwen2.5-coder

# Database
DATABASE_PATH=./data/app.db

# Server
SERVER_HOST=localhost
SERVER_PORT=8080
```

## 🎯 Supported AI Providers

| Provider | API Key Required | Local |
|----------|------------------|-------|
| **Gemini** | ✅ Free tier | ❌ |
| **Ollama** | ❌ | ✅ |
| **OpenAI** | ✅ Paid | ❌ |
| **Groq** | ✅ Free tier | ❌ |

## 📚 Educational Purpose

This project demonstrates:
- Go systems programming
- Wails desktop development
- Clean architecture patterns
- SQLite with migrations
- AI provider integration
- Repository pattern
- Event-driven design

## 🔧 Development

```bash
# Run in dev mode
wails dev

# Build for production
wails build

# Run Go tests
go test ./...
```

## 📄 License

MIT - For educational use.
