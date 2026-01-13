package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents an LLM API client
type Client struct {
	config *Config
	client *http.Client
}

// NewClient creates a new API client
func NewClient(cfg *Config) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.API.Timeout) * time.Second,
		},
	}
}

// ChatRequest represents an OpenAI-compatible chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents an OpenAI-compatible chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Translate translates text to the specified language
func (c *Client) Translate(text, targetLanguage string) (string, error) {
	prompt := c.buildPrompt(text, targetLanguage)

	request := ChatRequest{
		Model:       c.config.API.Model,
		Temperature: c.config.API.Temperature,
		MaxTokens:   c.config.API.MaxTokens,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.API.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.API.Key)

	if c.config.Advanced.Debug {
		fmt.Printf("[DEBUG] API Request: %s\n", string(jsonData))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if c.config.Advanced.Debug {
		fmt.Printf("[DEBUG] API Response: %s\n", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	translation := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	return translation, nil
}

// buildPrompt builds the translation prompt
func (c *Client) buildPrompt(text, targetLanguage string) string {
	template := c.config.Advanced.PromptTemplate
	if template == "" {
		template = `You are a professional translator. Translate the following text to {language}.
Only return the translated text without any explanation or additional content.

Text: {input_text}`
	}

	langName := getLanguageName(targetLanguage)
	prompt := strings.ReplaceAll(template, "{language}", langName)
	prompt = strings.ReplaceAll(prompt, "{input_text}", text)

	return prompt
}

// getLanguageName returns the full name of a language code
