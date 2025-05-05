package asst

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func init() {
	// ensure defaults for tests
	viper.SetDefault("api.base_url", "http://example.com")
	viper.SetDefault("api.key", "test-key")
	viper.SetDefault("api.model", "test-model")
	viper.SetDefault("assistant.system_msg_tmpl", "sysmsg: %s")
	viper.SetDefault("assistant.shell", "testshell")
	viper.SetDefault("api.timeout", 1)
}

func TestBuildRequest(t *testing.T) {
	prompt := "hello"
	req, err := buildRequest(prompt)
	if err != nil {
		t.Fatalf("buildRequest error: %v", err)
	}
	// check method
	if req.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", req.Method)
	}
	// check URL
	if req.URL.String() != "http://example.com/chat/completions" {
		t.Errorf("unexpected URL: %s", req.URL.String())
	}
	// check header
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type: application/json, got %s", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer test-key" {
		t.Errorf("expected Authorization: Bearer test-key, got %s", req.Header.Get("Authorization"))
	}
	// check body
	b, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("reading request body: %v", err)
	}
	var r ChatCompletionRequest
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if r.Model != "test-model" {
		t.Errorf("expected model test-model, got %s", r.Model)
	}
	if len(r.Messages) != 2 || r.Messages[1].Content != prompt {
		t.Errorf("unexpected messages: %+v", r.Messages)
	}
}

func TestSendRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("response body"))
		}),
	)
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	body, err := sendRequest(req)
	if err != nil {
		t.Fatalf("sendRequest error: %v", err)
	}
	if string(body) != "response body" {
		t.Errorf("expected 'response body', got %q", body)
	}
}

func TestSendRequest_Timeout(t *testing.T) {
	viper.SetDefault("api.timeout", 0.001)
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Millisecond)
		}),
	)
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	_, err := sendRequest(req)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		t.Fatalf("expected timeout error, got %v", err)
	}
}

func TestParseAssistantResponse(t *testing.T) {
	jsonStr := `{"choices":[{"message":{"content":"hi"}}]}`
	content, err := parseAssistantResponse([]byte(jsonStr))
	if err != nil {
		t.Fatalf("parseAssistantResponse error: %v", err)
	}
	if content != "hi" {
		t.Errorf("expected content hi, got %s", content)
	}
}

func TestExtractCodeBlock_Success(t *testing.T) {
	resp := "before ```lang\n code here\n``` after"
	code, err := extractCodeBlock(resp)
	if err != nil {
		t.Fatalf("extractCodeBlock error: %v", err)
	}
	if code != "code here" {
		t.Errorf("expected 'code here', got %q", code)
	}
}

func TestExtractCodeBlock_NoFence(t *testing.T) {
	_, err := extractCodeBlock("no fences here")
	if !errors.Is(err, ErrNoCodeBlock) {
		t.Errorf("expected ErrNoCodeBlock, got %v", err)
	}
}
