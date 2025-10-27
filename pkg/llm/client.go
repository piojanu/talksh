package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type requestBody struct {
	Model    string         `json:"model"`
	Messages []requestEntry `json:"messages"`
}

type requestEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseEnvelope struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func CallLLM(prompt, content string) (string, error) {
	baseURL := viper.GetString("api.base_url")
	key := viper.GetString("api.key")
	model := viper.GetString("api.model")
	timeout := viper.GetInt("api.timeout")
	promptPlaceholder := viper.GetString("llm.prompt_placeholder")
	systemMessage := viper.GetString("llm.system_message")

	if baseURL == "" || key == "" || model == "" {
		return "", fmt.Errorf("api.base_url, api.key, and api.model must be configured")
	}

	endpointURL := fmt.Sprintf("%s/chat/completions", baseURL)
	userMessage := strings.ReplaceAll(prompt, promptPlaceholder, content)

	payload := requestBody{
		Model: model,
		Messages: []requestEntry{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: userMessage},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", endpointURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call LLM: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("server error %d: %s", resp.StatusCode, respBody)
	}

	var envelope responseEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if envelope.Error != nil {
		return "", fmt.Errorf("llm error: %s", envelope.Error.Message)
	}
	if len(envelope.Choices) == 0 {
		return "", fmt.Errorf("llm returned no response")
	}
	return envelope.Choices[0].Message.Content, nil
}
