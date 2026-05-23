package handlers

import (
"ai-desktop-assistant/backend/internal/config"
"ai-desktop-assistant/backend/internal/services/ai"
"ai-desktop-assistant/backend/internal/services/storage"
)

type Handlers struct {
storage *storage.StorageService
ai      *ai.AIService
config  *config.Config
}

func New(storage *storage.StorageService, ai *ai.AIService, cfg *config.Config) *Handlers {
return &Handlers{
storage: storage,
ai:      ai,
config:  cfg,
}
}

func (h *Handlers) GetStorage() *storage.StorageService {
return h.storage
}

func (h *Handlers) GetAI() *ai.AIService {
return h.ai
}

func (h *Handlers) GetConfig() *config.Config {
return h.config
}
