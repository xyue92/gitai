package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient handles communication with local Ollama server
type OllamaClient struct {
	BaseURL      string
	Model        string
	Client       *http.Client
	EnableStream bool // Enable streaming output
}

// OllamaRequest represents the request structure for Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	Error     string `json:"error,omitempty"`
}

// NewOllamaClient creates a new Ollama client with timeout configuration
func NewOllamaClient(model string) *OllamaClient {
	return &OllamaClient{
		BaseURL: "http://localhost:11434",
		Model:   model,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate sends a prompt to Ollama and returns the generated text
func (c *OllamaClient) Generate(prompt string) (string, error) {
	// Check if Ollama is running
	if err := c.checkConnection(); err != nil {
		return "", fmt.Errorf("cannot connect to Ollama: %w\nPlease make sure Ollama is running:\n  $ ollama serve", err)
	}

	// Prepare request
	reqBody := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request
	url := c.BaseURL + "/api/generate"
	resp, err := c.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Try to extract error message
		var errResp OllamaResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			if contains(errResp.Error, "model") && contains(errResp.Error, "not found") {
				return "", fmt.Errorf("model '%s' not found\nInstall it with:\n  $ ollama pull %s", c.Model, c.Model)
			}
			return "", fmt.Errorf("ollama error: %s", errResp.Error)
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	return ollamaResp.Response, nil
}

// checkConnection checks if Ollama server is reachable
func (c *OllamaClient) checkConnection() error {
	url := c.BaseURL + "/api/tags"
	resp, err := c.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status code: %d", resp.StatusCode)
	}

	return nil
}

// GenerateStream sends a prompt to Ollama and streams the generated text
// onChunk is called for each chunk of text received
func (c *OllamaClient) GenerateStream(prompt string, onChunk func(chunk string)) (string, error) {
	// Check if Ollama is running
	if err := c.checkConnection(); err != nil {
		return "", fmt.Errorf("cannot connect to Ollama: %w\nPlease make sure Ollama is running:\n  $ ollama serve", err)
	}

	// Prepare request
	reqBody := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request
	url := c.BaseURL + "/api/generate"
	resp, err := c.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp OllamaResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			if contains(errResp.Error, "model") && contains(errResp.Error, "not found") {
				return "", fmt.Errorf("model '%s' not found\nInstall it with:\n  $ ollama pull %s", c.Model, c.Model)
			}
			return "", fmt.Errorf("ollama error: %s", errResp.Error)
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read streaming response
	decoder := json.NewDecoder(resp.Body)
	var fullResponse string

	for {
		var chunk OllamaResponse
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to decode chunk: %w", err)
		}

		if chunk.Error != "" {
			return "", fmt.Errorf("ollama error: %s", chunk.Error)
		}

		if chunk.Response != "" {
			fullResponse += chunk.Response
			if onChunk != nil {
				onChunk(chunk.Response)
			}
		}

		if chunk.Done {
			break
		}
	}

	return fullResponse, nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(bytes.Contains([]byte(s), []byte(substr))))
}
