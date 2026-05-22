package ai

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

type AIService struct {
apiKey string
model  string
client *http.Client
}

func NewGeminiService() *AIService {
return &AIService{
model: "gemini-2.0-flash",
client: &http.Client{Timeout: 60 * time.Second},
}
}

func (a *AIService) SetAPIKey(key string) {
a.apiKey = key
}

func (a *AIService) SetModel(model string) {
a.model = model
}

func (a *AIService) GetModel() string {
return a.model
}

func (a *AIService) Chat(ctx context.Context, message string) (string, error) {
if a.apiKey == "" {
return "", fmt.Errorf("API key not configured")
}

url := fmt.Sprintf(
"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
a.model, a.apiKey,
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

resp, err := a.client.Do(req)
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

// GetAvailableModels returns available Gemini models
func (a *AIService) GetAvailableModels() []string {
return []string{
"gemini-2.0-flash",
"gemini-1.5-flash", 
"gemini-1.5-pro",
}
}
