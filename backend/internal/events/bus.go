package events

import (
"sync"
)

// Event type identifiers
const (
EventChatMessage     = "chat:message"
EventChatResponse    = "chat:response"
EventOverlayToggle   = "overlay:toggle"
EventOverlayOpacity  = "overlay:opacity"
EventScreenshot      = "screenshot:captured"
EventOCRComplete     = "ocr:complete"
EventAudioCapture    = "audio:captured"
Event Transcription  = "transcription:complete"
EventSettingsChanged = "settings:changed"
EventSessionCreated  = "session:created"
EventError           = "error"
)

// EventHandler is a function that handles events
type EventHandler func(payload interface{})

// Bus is the event bus for publish/subscribe pattern
type Bus struct {
handlers map[string][]EventHandler
mu       sync.RWMutex
}

// New creates a new event bus
func New() *Bus {
return &Bus{
handlers: make(map[string][]EventHandler),
}
}

// Subscribe registers a handler for an event type
func (b *Bus) Subscribe(eventType string, handler EventHandler) {
b.mu.Lock()
defer b.mu.Unlock()

b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Unsubscribe removes a handler for an event type
func (b *Bus) Unsubscribe(eventType string, handler EventHandler) {
b.mu.Lock()
defer b.mu.Unlock()

handlers := b.handlers[eventType]
for i, h := range handlers {
if &h == &handler {
b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
return
}
}
}

// Publish publishes an event to all registered handlers
func (b *Bus) Publish(eventType string, payload interface{}) {
b.mu.RLock()
handlers := b.handlers[eventType]
b.mu.RUnlock()

for _, handler := range handlers {
go handler(payload)
}
}

// Broadcast publishes to multiple event types
func (b *Bus) Broadcast(eventTypes []string, payload interface{}) {
for _, eventType := range eventTypes {
b.Publish(eventType, payload)
}
}

// Clear removes all handlers
func (b *Bus) Clear() {
b.mu.Lock()
defer b.mu.Unlock()
b.handlers = make(map[string][]EventHandler)
}

// HandlerCount returns the number of handlers for an event type
func (b *Bus) HandlerCount(eventType string) int {
b.mu.RLock()
defer b.mu.RUnlock()
return len(b.handlers[eventType])
}
