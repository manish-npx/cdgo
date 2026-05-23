import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

// Wails runtime bridge - this provides the window.go object
declare global {
  interface Window {
    go: {
      main: {
        GetAPIKey(): Promise<string>
        SetAPIKey(key: string): Promise<void>
        GetAIModel(): Promise<string>
        SetAIModel(model: string): Promise<void>
        GetSettings(): Promise<Record<string, unknown>>
        SaveSettings(settings: Record<string, unknown>): Promise<void>
        SendMessage(message: string): Promise<string>
        GetChatHistory(): Promise<Array<Record<string, unknown>>>
      }
    }
  }
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
