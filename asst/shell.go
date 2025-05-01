/*
Copyright 2025 Piotr Januszewski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package asst

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

var ErrNoCodeBlock = errors.New("no code-block found")

func SuggestCommand(prompt string) (string, error) {
	req, err := buildRequest(prompt)
	if err != nil {
		return "", err
	}
	body, err := sendRequest(req)
	if err != nil {
		return "", err
	}
	assistantResponse, err := parseAssistantResponse(body)
	if err != nil {
		return "", err
	}
	return extractCodeBlock(assistantResponse)
}

func buildRequest(prompt string) (*http.Request, error) {
	requestURL := fmt.Sprintf("%s/chat/completions", viper.GetString("api.base_url"))
	modelName := viper.GetString("api.model")
	messages := []Message{
		{
			Role: "system",
			Content: fmt.Sprintf(
				viper.GetString("assistant.system_msg_tmpl"),
				viper.GetString("assistant.shell"),
			),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
	b, err := json.Marshal(ChatCompletionRequest{
		Model:    modelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal messages: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func sendRequest(req *http.Request) ([]byte, error) {
	// allow fractional timeouts (mostly for tests)
	timeout := time.Duration(
		viper.GetFloat64("api.timeout") * float64(time.Second))
	client := http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to call API: status code %d: %q", resp.StatusCode, body,
		)
	}
	return body, nil
}

func parseAssistantResponse(body []byte) (string, error) {
	var partial struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &partial); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return partial.Choices[0].Message.Content, nil
}

func extractCodeBlock(assistantResponse string) (string, error) {
	parts := strings.Split(assistantResponse, "```")
	if len(parts) < 3 {
		return "", fmt.Errorf(
			"failed to parse the assistant's response: %q: %w",
			assistantResponse,
			ErrNoCodeBlock,
		)
	}
	codeBlock := parts[1]
	// Drop the optional language tag
	if nl := strings.IndexRune(codeBlock, '\n'); nl >= 0 {
		codeBlock = codeBlock[nl+1:]
	}
	return strings.TrimSpace(codeBlock), nil
}
