import { useState, useEffect } from 'react'
import { Settings, Send, Trash2, Bot, User } from 'lucide-react'

interface Message {
  id: string
  role: 'user' | 'ai'
  content: string
  timestamp: Date
}

function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [apiKey, setApiKey] = useState('')
  const [aiModel, setAiModel] = useState('gemini-2.0-flash')
  const [showSettings, setShowSettings] = useState(false)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadSettings()
  }, [])

  const loadSettings = async () => {
    try {
      if (window.go?.main) {
        const key = await window.go.main.GetAPIKey()
        const model = await window.go.main.GetAIModel()
        setApiKey(key || '')
        setAiModel(model || 'gemini-2.0-flash')
        
        if (key) {
          const history = await window.go.main.GetChatHistory()
          if (history && history.length > 0) {
            const loaded = history.map((h: any) => ({
              id: h.id || Date.now().toString(),
              role: 'user' as const,
              content: h.userMsg || '',
              timestamp: new Date(h.timestamp || Date.now())
            }))
            setMessages(loaded)
          }
        }
      }
    } catch (e) {
      // Wails not available in dev mode
      console.log('Running in development mode')
    }
  }

  const saveApiKey = async () => {
    try {
      if (window.go?.main) {
        await window.go.main.SetAPIKey(apiKey)
        await window.go.main.SetAIModel(aiModel)
      }
    } catch (e) {
      console.log('Could not save API key')
    }
    setShowSettings(false)
  }

  const sendMessage = async () => {
    if (!input.trim()) return
    if (!apiKey) {
      setShowSettings(true)
      return
    }

    const userMsg = input
    setInput('')
    setLoading(true)

    // Add user message
    setMessages(prev => [...prev, {
      id: Date.now().toString(),
      role: 'user',
      content: userMsg,
      timestamp: new Date()
    }])

    try {
      let response = ''
      
      if (window.go?.main) {
        // Use Wails backend
        response = await window.go.main.SendMessage(userMsg)
      } else {
        // Development fallback - call API directly
        response = await callGeminiAPI(userMsg, apiKey, aiModel)
      }

      setMessages(prev => [...prev, {
        id: (Date.now() + 1).toString(),
        role: 'ai',
        content: response,
        timestamp: new Date()
      }])
    } catch (e: any) {
      setMessages(prev => [...prev, {
        id: (Date.now() + 1).toString(),
        role: 'ai',
        content: `Error: ${e.message}`,
        timestamp: new Date()
      }])
    }

    setLoading(false)
  }

  const callGeminiAPI = async (message: string, key: string, model: string) => {
    const response = await fetch(
      `https://generativelanguage.googleapis.com/v1beta/models/${model}:generateContent?key=${key}`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          contents: [{ parts: [{ text: message }] }]
        })
      }
    )
    
    const data = await response.json()
    if (!response.ok) throw new Error(data.error?.message || 'API Error')
    
    return data.candidates?.[0]?.content?.parts?.[0]?.text || 'No response'
  }

  const clearChat = () => setMessages([])

  return (
    <div className="app">
      {/* Header */}
      <div className="header">
        <div className="title">
          <Bot size={24} />
          <span>AI Desktop Assistant</span>
        </div>
        <button className="icon-btn" onClick={() => setShowSettings(!showSettings)}>
          <Settings size={20} />
        </button>
      </div>

      {/* Settings Panel */}
      {showSettings && (
        <div className="settings-panel">
          <h3>⚙️ Settings</h3>
          
          <div className="setting-group">
            <label>Gemini API Key</label>
            <input
              type="password"
              value={apiKey}
              onChange={e => setApiKey(e.target.value)}
              placeholder="Enter your API key (get free at aistudio.google.com)"
            />
          </div>

          <div className="setting-group">
            <label>AI Model</label>
            <select value={aiModel} onChange={e => setAiModel(e.target.value)}>
              <option value="gemini-2.0-flash">Gemini 2.0 Flash (Fast)</option>
              <option value="gemini-1.5-flash">Gemini 1.5 Flash</option>
              <option value="gemini-1.5-pro">Gemini 1.5 Pro</option>
            </select>
          </div>

          <button className="btn-primary" onClick={saveApiKey}>
            Save Settings
          </button>
        </div>
      )}

      {/* Messages */}
      <div className="messages">
        {messages.length === 0 && !apiKey && (
          <div className="welcome">
            <Bot size={48} />
            <h2>Welcome to AI Desktop Assistant</h2>
            <p>Enter your Gemini API key in settings to start chatting</p>
            <button className="btn-primary" onClick={() => setShowSettings(true)}>
              Setup API Key
            </button>
          </div>
        )}

        {messages.map(msg => (
          <div key={msg.id} className={`message ${msg.role}`}>
            <div className="avatar">
              {msg.role === 'user' ? <User size={16} /> : <Bot size={16} />}
            </div>
            <div className="content">{msg.content}</div>
          </div>
        ))}

        {loading && (
          <div className="message ai">
            <div className="avatar"><Bot size={16} /></div>
            <div className="content typing">Thinking...</div>
          </div>
        )}
      </div>

      {/* Input */}
      <div className="input-area">
        {messages.length > 0 && (
          <button className="icon-btn" onClick={clearChat} title="Clear chat">
            <Trash2 size={18} />
          </button>
        )}
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && sendMessage()}
          placeholder={apiKey ? "Type your message..." : "Enter API key first"}
          disabled={loading}
        />
        <button className="icon-btn send" onClick={sendMessage} disabled={loading || !input.trim()}>
          <Send size={18} />
        </button>
      </div>

      <style>{`
        .app {
          display: flex;
          flex-direction: column;
          height: 100vh;
          background: #0f172a;
        }
        .header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px 16px;
          background: #1e293b;
          border-bottom: 1px solid #334155;
        }
        .title {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 16px;
          font-weight: 600;
          color: #667eea;
        }
        .icon-btn {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 36px;
          height: 36px;
          border: none;
          border-radius: 8px;
          background: #334155;
          color: #94a3b8;
          cursor: pointer;
          transition: all 0.2s;
        }
        .icon-btn:hover {
          background: #475569;
          color: #e2e8f0;
        }
        .icon-btn.send {
          background: #667eea;
          color: white;
        }
        .icon-btn.send:hover {
          background: #764ba2;
        }
        .settings-panel {
          padding: 16px;
          background: #1e293b;
          border-bottom: 1px solid #334155;
        }
        .settings-panel h3 {
          margin-bottom: 16px;
          color: #94a3b8;
        }
        .setting-group {
          margin-bottom: 12px;
        }
        .setting-group label {
          display: block;
          margin-bottom: 6px;
          font-size: 14px;
          color: #94a3b8;
        }
        .setting-group input,
        .setting-group select {
          width: 100%;
          padding: 10px 12px;
          border: 1px solid #334155;
          border-radius: 8px;
          background: #0f172a;
          color: #e2e8f0;
          font-size: 14px;
        }
        .setting-group input:focus,
        .setting-group select:focus {
          outline: none;
          border-color: #667eea;
        }
        .btn-primary {
          width: 100%;
          padding: 12px;
          border: none;
          border-radius: 8px;
          background: linear-gradient(135deg, #667eea, #764ba2);
          color: white;
          font-size: 14px;
          font-weight: 500;
          cursor: pointer;
          transition: transform 0.2s;
        }
        .btn-primary:hover {
          transform: scale(1.02);
        }
        .messages {
          flex: 1;
          overflow-y: auto;
          padding: 16px;
        }
        .welcome {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          height: 100%;
          text-align: center;
          color: #94a3b8;
        }
        .welcome h2 {
          margin: 16px 0 8px;
          color: #e2e8f0;
        }
        .welcome p {
          margin-bottom: 24px;
        }
        .message {
          display: flex;
          gap: 12px;
          margin-bottom: 16px;
        }
        .message.user {
          flex-direction: row-reverse;
        }
        .avatar {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 32px;
          height: 32px;
          border-radius: 50%;
          background: #334155;
          color: #94a3b8;
          flex-shrink: 0;
        }
        .message.user .avatar {
          background: #667eea;
          color: white;
        }
        .message.ai .avatar {
          background: #475569;
        }
        .content {
          max-width: 70%;
          padding: 12px 16px;
          border-radius: 12px;
          background: #1e293b;
          color: #e2e8f0;
          line-height: 1.5;
          white-space: pre-wrap;
        }
        .message.user .content {
          background: #334155;
        }
        .typing {
          color: #94a3b8;
          font-style: italic;
        }
        .input-area {
          display: flex;
          gap: 8px;
          padding: 16px;
          background: #1e293b;
          border-top: 1px solid #334155;
        }
        .input-area input {
          flex: 1;
          padding: 12px 16px;
          border: 1px solid #334155;
          border-radius: 24px;
          background: #0f172a;
          color: #e2e8f0;
          font-size: 14px;
        }
        .input-area input:focus {
          outline: none;
          border-color: #667eea;
        }
      `}</style>
    </div>
  )
}

export default App
