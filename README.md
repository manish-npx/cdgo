# AI Desktop Assistant 🤖

A production-grade desktop AI assistant built with **Go + Wails + React** for educational purposes.

## Features

- **Standalone Desktop App** - No browser needed, runs natively on Windows/Mac/Linux
- **Gemini AI Integration** - Uses Google Gemini API for AI chat
- **Dark Theme UI** - Modern, sleek interface
- **Settings Panel** - Enter your own API key
- **Session History** - Chat history saved locally
- **Always On Top** - Keeps window visible while working

## 🚀 Quick Start

### Prerequisites

1. **Go 1.21+** - [Download](https://go.dev/dl/)
2. **Node.js 18+** - [Download](https://nodejs.org/)
3. **Wails CLI** - Install via: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Setup Steps

```bash
# 1. Clone the repository
git clone https://github.com/manish-npx/cdgo.git
cd cdgo/ai-desktop-assistant

# 2. Install frontend dependencies
cd frontend
npm install

# 3. Build frontend
npm run build

# 4. Go back and run with Wails
cd ..
wails dev
```

### First Run

1. The app opens with a **Settings panel**
2. Enter your **Gemini API key** (get free from [aistudio.google.com](https://aistudio.google.com/app/apikey))
3. Select your preferred **AI model**
4. Click **Save Settings**
5. Start chatting!

## 📁 Project Structure

```
ai-desktop-assistant/
├── backend/                    # Go + Wails backend
│   ├── cmd/main.go            # Wails app entry
│   └── internal/
│       ├── config/            # Configuration management
│       ├── handlers/          # Wails IPC handlers
│       └── services/
│           ├── ai/           # Gemini AI service
│           └── storage/      # Session storage
├── frontend/                   # React frontend
│   ├── src/
│   │   ├── App.tsx           # Main React component
│   │   └── main.tsx          # Entry point
│   └── package.json
└── README.md
```

## 🔑 Getting Gemini API Key

1. Go to [https://aistudio.google.com/app/apikey](https://aistudio.google.com/app/apikey)
2. Click "Create API Key"
3. Copy the key
4. Paste it in the app's Settings panel

**Free Tier Available:**
- Gemini 2.0 Flash: 15 requests/min
- Gemini 1.5 Flash: 15 requests/min
- Gemini 1.5 Pro: Paid

## 🎯 How It Works

```
┌──────────────────────────┐
│     Desktop Window        │
│  ┌────────────────────┐  │
│  │   React UI         │  │
│  │   (TypeScript)     │  │
│  └────────┬───────────┘  │
│           │ Wails IPC     │
│  ┌────────▼───────────┐  │
│  │   Go Backend        │  │
│  │   - AI Service      │  │
│  │   - Storage         │  │
│  │   - Config          │  │
│  └────────┬───────────┘  │
│           │              │
│  ┌────────▼───────────┐  │
│  │   Gemini API        │  │
│  │   (Your API Key)    │  │
│  └────────────────────┘  │
└──────────────────────────┘
```

## ⚙️ Configuration

Settings are stored at:
- **Windows:** `C:\Users\<You>\.ai-desktop-assistant\config.json`
- **Linux/Mac:** `~/.ai-desktop-assistant/config.json`

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| Desktop Framework | Wails v2 |
| Backend | Go |
| Frontend | React 18 + TypeScript |
| State | React hooks |
| AI | Google Gemini API |

## 📚 Educational Purpose

This project demonstrates:
- Go systems programming
- Wails desktop development
- React frontend integration
- Clean architecture patterns
- API integration patterns
- Desktop app deployment

## 🔧 Development

```bash
# Run in dev mode with hot reload
wails dev

# Build for production
wails build

# Run backend only
go run backend/cmd/main.go
```

## 📄 License

MIT - For educational purposes.
