package ai

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed prompts/*.json
var promptFS embed.FS

const promptsDir = "prompts"

// PromptDefinition is a versioned prompt document loaded from JSON.
type PromptDefinition struct {
	Version           string          `json:"version"`
	ID                string          `json:"id"`
	Model             string          `json:"model"`
	Temperature       float64         `json:"temperature"`
	MaxOutputTokens   int             `json:"max_output_tokens"`
	SystemInstruction string          `json:"system_instruction"`
	UserInputShape    json.RawMessage `json:"user_input_shape"`
	OutputSchema      map[string]any  `json:"output_schema"`
	SchemaName        string          `json:"schema_name"`
	ExampleInput      json.RawMessage `json:"example_input"`
	ExampleOutput     json.RawMessage `json:"example_output"`
}

// Registry loads and caches embedded prompt definitions.
type Registry struct {
	mu      sync.RWMutex
	byID    map[string]PromptDefinition
	byFile  map[string]PromptDefinition
}

func NewRegistry() (*Registry, error) {
	r := &Registry{
		byID:   make(map[string]PromptDefinition),
		byFile: make(map[string]PromptDefinition),
	}
	entries, err := promptFS.ReadDir(promptsDir)
	if err != nil {
		return nil, fmt.Errorf("read prompts dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := promptFS.ReadFile(promptsDir + "/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read prompt %s: %w", entry.Name(), err)
		}
		var def PromptDefinition
		if err := json.Unmarshal(data, &def); err != nil {
			return nil, fmt.Errorf("parse prompt %s: %w", entry.Name(), err)
		}
		if def.ID == "" {
			return nil, fmt.Errorf("prompt %s missing id", entry.Name())
		}
		if def.Version == "" {
			return nil, fmt.Errorf("prompt %s missing version", entry.Name())
		}
		if def.SchemaName == "" {
			def.SchemaName = def.ID
		}
		r.byID[def.ID] = def
		r.byFile[entry.Name()] = def
	}
	return r, nil
}

func (r *Registry) Get(id string) (PromptDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	def, ok := r.byID[id]
	if !ok {
		return PromptDefinition{}, fmt.Errorf("prompt not found: %s", id)
	}
	return def, nil
}

// VersionedName returns the canonical filename-style version label, e.g. baseline_diagnosis.v1.
func (d PromptDefinition) VersionedName() string {
	return d.ID + "." + d.Version
}
