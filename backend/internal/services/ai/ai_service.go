package ai

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"time"

"ai-desktop-assistant/backend/internal/config"
)

// AIService handles AI provider communication
type AIService struct {
provider   string
geminiKey  string
geminiModel string
ollamaHost string
ollamaModel string
client     *http.Client
}

// NewService creates a new AI service
func NewService(cfg *config.Config) *AIService {
return &AIService{
provider:    cfg.AI.Provider,
geminiKey:   cfg.AI.GeminiAPIKey,
geminiModel: cfg.AI.GeminiModel,
ollamaHost:  cfg.AI.OllamaHost,
ollamaModel: cfg.AI.OllamaModel,
client: &http.Client{
Timeout: 120 * time.Second,
},
}
}

func (s *AIService) SetProvider(provider string) {
s.provider = provider
}

func (s *AIService) SetAPIKey(key string) {
s.geminiKey = key
}

func (s *AIService) SetModel(model string) {
s.geminiModel = model
}

func (s *AIService) SetOllamaHost(host string) {
s.ollamaHost = host
}

func (s *AIService) SetOllamaModel(model string) {
s.ollamaModel = model
}

// Chat sends a message to the configured AI provider
func (s *AIService) Chat(ctx context.Context, message string) (string, error) {
switch s.provider {
case "ollama":
return s.chatOllama(ctx, message)
case "gemini":
return s.chatGemini(ctx, message)
default:
if s.geminiKey != "" {
return s.chatGemini(ctx, message)
}
return "", fmt.Errorf("no AI provider configured")
}
}

// chatGemini calls Google Gemini API
func (s *AIService) chatGemini(ctx context.Context, message string) (string, error) {
if s.geminiKey == "" {
return "", fmt.Errorf("Gemini API key not set")
}

url := fmt.Sprintf(
"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
s.geminiModel, s.geminiKey,
)

reqBody := map[string]interface{}{
"contents": []map[string]interface{}{
{"parts": []map[string]string{{"text": message}}},
},
}

jsonData, err := json.Marshal(reqBody)
if err != nil {
return "", fmt.Errorf("failed to marshal: %w", err)
}

req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
if err != nil {
return "", fmt.Errorf("failed to create request: %w", err)
}
req.Header.Set("Content-Type", "application/json")

resp, err := s.client.Do(req)
if err != nil {
return "", fmt.Errorf("request failed: %w", err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return "", fmt.Errorf("failed to read response: %w", err)
}

if resp.StatusCode != http.StatusOK {
return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
}

var result map[string]interface{}
if err := json.Unmarshal(body, &result); err != nil {
return "", fmt.Errorf("failed to decode: %w", err)
}

candidates, ok := result["candidates"].([]interface{})
if !ok || len(candidates) == 0 {
return "", fmt.Errorf("no response from API")
}

candidate := candidates[0].(map[string]interface{})
content := candidate["content"].(map[string]interface{})
parts := content["parts"].([]interface{})
part := parts[0].(map[string]interface{})

text, ok := part["text"].(string)
if !ok {
return "", fmt.Errorf("no text in response")
}

return text, nil
}

// chatOllama calls local Ollama server
func (s *AIService) chatOllama(ctx context.Context, message string) (string, error) {
url := fmt.Sprintf("%s/api/chat", s.ollamaHost)

reqBody := map[string]interface{}{
"model":  s.ollamaModel,
"stream": false,
"messages": []map[string]string{
{"role": "user", "content": message},
},
}

jsonData, err := json.Marshal(reqBody)
if err != nil {
return "", fmt.Errorf("failed to marshal: %w", err)
}

req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
if err != nil {
return "", fmt.Errorf("failed to create request: %w", err)
}
req.Header.Set("Content-Type", "application/json")

resp, err := s.client.Do(req)
if err != nil {
return "", fmt.Errorf("request failed: %w", err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
return "", fmt.Errorf("failed to read response: %w", err)
}

if resp.StatusCode != http.StatusOK {
return "", fmt.Errorf("Ollama error (%d): %s", resp.StatusCode, string(body))
}

var result struct {
Message struct {
Content string `json:"content"`
} `json:"message"`
}

if err := json.Unmarshal(body, &result); err != nil {
return "", fmt.Errorf("failed to decode: %w", err)
}

return result.Message.Content, nil
}

// GetAvailableModels returns models for each provider
func (s *AIService) GetAvailableModels() []string {
return []string{
"gemini-2.0-flash",
"gemini-1.5-flash",
"gemini-1.5-pro",
"llama3.2",
"qwen2.5-coder",
}
}
