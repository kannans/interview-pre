package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.openai.com/v1"

// ClientConfig holds OpenAI Responses API settings.
type ClientConfig struct {
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
}

// CallMeta captures observability fields for a single AI request.
type CallMeta struct {
	PromptID      string
	PromptVersion string
	Model         string
	Latency       time.Duration
	InputTokens   int
	OutputTokens  int
	TotalTokens   int
	ResponseID    string
}

// StructuredRequest is sent to POST /v1/responses.
type StructuredRequest struct {
	Model            string
	Instructions     string
	UserInput        string
	SchemaName       string
	OutputSchema     map[string]any
	Temperature      float64
	MaxOutputTokens  int
	PromptID         string
	PromptVersion    string
}

type responsesAPIRequest struct {
	Model           string         `json:"model"`
	Instructions    string         `json:"instructions,omitempty"`
	Input           string         `json:"input"`
	Temperature     float64        `json:"temperature,omitempty"`
	MaxOutputTokens int            `json:"max_output_tokens,omitempty"`
	Text            responsesText  `json:"text,omitempty"`
}

type responsesText struct {
	Format responsesFormat `json:"format"`
}

type responsesFormat struct {
	Type   string         `json:"type"`
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
	Strict bool           `json:"strict"`
}

type responsesAPIResponse struct {
	ID         string `json:"id"`
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Client wraps the OpenAI Responses API.
type Client struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg ClientConfig) *Client {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 90 * time.Second}
	}
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &Client{
		apiKey:     cfg.APIKey,
		model:      model,
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// CompleteStructured sends a structured-output request and returns raw JSON bytes.
func (c *Client) CompleteStructured(ctx context.Context, req StructuredRequest) (json.RawMessage, CallMeta, error) {
	if c.apiKey == "" {
		return nil, CallMeta{}, ErrNotConfigured
	}

	model := req.Model
	if model == "" {
		model = c.model
	}

	body := responsesAPIRequest{
		Model:           model,
		Instructions:    req.Instructions,
		Input:           req.UserInput,
		Temperature:     req.Temperature,
		MaxOutputTokens: req.MaxOutputTokens,
		Text: responsesText{
			Format: responsesFormat{
				Type:   "json_schema",
				Name:   req.SchemaName,
				Schema: req.OutputSchema,
				Strict: true,
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, CallMeta{}, fmt.Errorf("marshal openai request: %w", err)
	}

	start := time.Now()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(payload))
	if err != nil {
		return nil, CallMeta{}, fmt.Errorf("create openai request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, CallMeta{}, fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, CallMeta{}, fmt.Errorf("read openai response: %w", err)
	}

	latency := time.Since(start)
	meta := CallMeta{
		PromptID:      req.PromptID,
		PromptVersion: req.PromptVersion,
		Model:         model,
		Latency:       latency,
	}

	if resp.StatusCode >= 400 {
		slog.Error("openai api error",
			"prompt_id", req.PromptID,
			"prompt_version", req.PromptVersion,
			"model", model,
			"latency_ms", latency.Milliseconds(),
			"status", resp.StatusCode,
			"body", string(respBody),
		)
		return nil, meta, fmt.Errorf("openai api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed responsesAPIResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, meta, fmt.Errorf("decode openai response envelope: %w", err)
	}

	meta.ResponseID = parsed.ID
	meta.InputTokens = parsed.Usage.InputTokens
	meta.OutputTokens = parsed.Usage.OutputTokens
	meta.TotalTokens = parsed.Usage.TotalTokens

	slog.Info("openai call completed",
		"prompt_id", req.PromptID,
		"prompt_version", req.PromptVersion,
		"model", model,
		"latency_ms", latency.Milliseconds(),
		"input_tokens", meta.InputTokens,
		"output_tokens", meta.OutputTokens,
		"total_tokens", meta.TotalTokens,
		"response_id", meta.ResponseID,
	)

	text := parsed.OutputText
	if text == "" {
		text = extractOutputText(parsed.Output)
	}
	if text == "" {
		if parsed.Error != nil {
			return nil, meta, fmt.Errorf("openai error: %s", parsed.Error.Message)
		}
		return nil, meta, ErrEmptyResponse
	}

	raw := json.RawMessage(text)
	if !json.Valid(raw) {
		return nil, meta, ErrInvalidJSON
	}

	return raw, meta, nil
}

func extractOutputText(output []struct {
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}) string {
	for _, item := range output {
		for _, part := range item.Content {
			if part.Text != "" {
				return part.Text
			}
		}
	}
	return ""
}
