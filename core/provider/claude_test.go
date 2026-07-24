package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Claude adapter unit tests using httptest.Server
// ---------------------------------------------------------------------------

func TestClaudeComplete(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		response: `{
			"id": "msg_test",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Hello from Claude!"}],
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 10, "output_tokens": 5}
		}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	resp, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if resp.Message.Content != "Hello from Claude!" {
		t.Errorf("content: got=%s", resp.Message.Content)
	}
	if resp.FinishReason != "end_turn" {
		t.Errorf("reason: got=%s", resp.FinishReason)
	}
	if resp.Usage.InputTokens != 10 || resp.Usage.OutputTokens != 5 {
		t.Errorf("usage: got=%+v", resp.Usage)
	}
}

func TestClaudeCompleteWithOptions(t *testing.T) {
	var captured struct {
		Model     string        `json:"model"`
		MaxTokens int           `json:"max_tokens"`
		System    string        `json:"system"`
		Messages  []interface{} `json:"messages"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"msg1","type":"message","role":"assistant","content":[{"type":"text","text":"ok"}],"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`))
	}))
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	_, err := cl.Complete(context.Background(), CompletionRequest{
		Model:       "claude-3-opus-20240229",
		Messages:    []Message{{Role: RoleUser, Content: "Test"}},
		System:      "Be concise.",
		MaxTokens:   500,
		Temperature: 0.3,
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}

	if captured.Model != "claude-3-opus-20240229" {
		t.Errorf("model: got=%s", captured.Model)
	}
	if captured.MaxTokens != 500 {
		t.Errorf("maxTokens: got=%d", captured.MaxTokens)
	}
	if captured.System != "Be concise." {
		t.Errorf("system: got=%s", captured.System)
	}
	if len(captured.Messages) != 1 {
		t.Errorf("messages: got=%d", len(captured.Messages))
	}
}

func TestClaudeCompleteMultiContent(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		response: `{
			"id": "msg_multi",
			"type": "message",
			"role": "assistant",
			"content": [
				{"type": "text", "text": "First part. "},
				{"type": "text", "text": "Second part."}
			],
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 5, "output_tokens": 10}
		}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	resp, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Write two sentences"}},
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	expected := "First part. Second part."
	if resp.Message.Content != expected {
		t.Errorf("content: got=%s want=%s", resp.Message.Content, expected)
	}
}

func TestClaudeCompleteEmptyResponse(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		response: `{
			"id": "msg_empty",
			"type": "message",
			"role": "assistant",
			"content": [],
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 5, "output_tokens": 0}
		}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	resp, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Say nothing"}},
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if resp.Message.Content != "" {
		t.Errorf("content: got=%s", resp.Message.Content)
	}
}

func TestClaudeStream(t *testing.T) {
	events := []string{
		`event: message_start
data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","content":[],"stop_reason":null,"usage":{"input_tokens":5,"output_tokens":0}}}`,
		``,
		`event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`,
		``,
		`event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" world"}}`,
		``,
		`event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":3}}`,
		``,
		`event: message_stop
data: {"type":"message_stop"}`,
	}

	srv := newClaudeStreamServer(t, events)
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	ch, err := cl.Stream(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}

	var full string
	chunkCount := 0
	for chunk := range ch {
		if chunk.Err != nil {
			t.Fatalf("chunk err: %v", chunk.Err)
		}
		full += chunk.Content
		chunkCount++
		if chunk.Done {
			if chunk.Usage.OutputTokens != 3 {
				t.Errorf("usage output: got=%d", chunk.Usage.OutputTokens)
			}
		}
	}

	if full != "Hello world" {
		t.Errorf("streamed content: got=%s", full)
	}
	if chunkCount < 3 {
		t.Errorf("too few chunks: %d", chunkCount)
	}
}

func TestClaudeStreamEmpty(t *testing.T) {
	events := []string{
		`event: message_start
data: {"type":"message_start","message":{"id":"msg_e","type":"message","role":"assistant","content":[],"stop_reason":null,"usage":{"input_tokens":1,"output_tokens":0}}}`,
		``,
		`event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":0}}`,
		``,
		`event: message_stop
data: {"type":"message_stop"}`,
	}

	srv := newClaudeStreamServer(t, events)
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	ch, err := cl.Stream(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Be silent"}},
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}

	var full string
	for chunk := range ch {
		if chunk.Err != nil {
			t.Fatalf("chunk err: %v", chunk.Err)
		}
		full += chunk.Content
	}
	if full != "" {
		t.Errorf("expected empty, got=%s", full)
	}
}

func TestClaudeRateLimited(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusTooManyRequests,
		response:   `{"type":"error","error":{"type":"rate_limit_error","message":"Rate limited"}}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	_, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	if !strings.Contains(err.Error(), ErrRateLimited.Error()) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeAuthFailed(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusUnauthorized,
		response:   `{"type":"error","error":{"type":"authentication_error","message":"Invalid API key"}}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	_, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected auth error")
	}
	if !strings.Contains(err.Error(), ErrAuthFailed.Error()) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeContextWindowExceeded(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusBadRequest,
		response:   `{"type":"error","error":{"type":"invalid_request_error","message":"prompt is too long: 200001 > 200000"}}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	_, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: strings.Repeat("a", 200001)}},
	})
	if err == nil {
		t.Fatal("expected context window error")
	}
	if !strings.Contains(err.Error(), ErrContextWindowExceeded.Error()) {
		t.Errorf("unexpected error: %v", err)
	}

	// Also test the pattern match via the response body path (non-JSON error)
	srv2 := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusBadRequest,
		response:   `prompt too long: context length exceeded`,
	})
	defer srv2.Close()

	cl2 := newTestClaude(t, srv2.URL)
	_, err = cl2.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), ErrContextWindowExceeded.Error()) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeProviderUnavailable(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusServiceUnavailable,
		response:   `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	_, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected unavailable error")
	}
	if !strings.Contains(err.Error(), ErrProviderUnavailable.Error()) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeBadRequest(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusBadRequest,
		response:   `{"type":"error","error":{"type":"invalid_request_error","message":"invalid model"}}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	_, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected bad request error")
	}
	if !strings.Contains(err.Error(), ErrBadRequest.Error()) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeServerError(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		statusCode: http.StatusInternalServerError,
		response:   `Internal Server Error`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	_, err := cl.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), ErrProviderUnavailable.Error()) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClaudeTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cl := NewClaude(
		WithClaudeBaseURL(srv.URL),
		WithClaudeAPIKey("test-key"),
		WithClaudeTimeout(5*time.Millisecond),
	)

	ctx := context.Background()
	_, err := cl.Complete(ctx, CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	t.Logf("timeout error (expected): %v", err)
}

func TestClaudeCapabilities(t *testing.T) {
	cl := NewClaude(WithClaudeAPIKey("test-key"))
	caps := cl.Capabilities()
	if caps.Provider != "anthropic" {
		t.Errorf("provider: got=%s", caps.Provider)
	}
	if !caps.Streaming {
		t.Error("should support streaming")
	}
	if len(caps.Models) == 0 {
		t.Error("should have models")
	}
}

func TestClaudeContextCancellation(t *testing.T) {
	srv := newClaudeTestServer(t, claudeTestHandler{
		response: `{
			"id": "msg_cc",
			"type": "message",
			"role": "assistant",
			"content": [{"type":"text","text":"should not get here"}],
			"stop_reason": "end_turn",
			"usage": {"input_tokens":1,"output_tokens":1}
		}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := cl.Complete(ctx, CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	t.Logf("cancellation error (expected): %v", err)
}

func TestClaudeStreamContextCancellation(t *testing.T) {
	events := []string{
		`event: message_start
data: {"type":"message_start","message":{"id":"msg_sc","type":"message","role":"assistant","content":[],"stop_reason":null,"usage":{"input_tokens":1,"output_tokens":0}}}`,
		``,
		`event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`,
	}

	srv := newClaudeStreamServer(t, events)
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := cl.Stream(ctx, CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}

	// Read a chunk, then cancel
	gotContent := false
	for chunk := range ch {
		if chunk.Err != nil {
			break
		}
		if chunk.Content != "" {
			gotContent = true
			cancel() // cancel mid-stream
		}
	}
	if !gotContent {
		t.Log("note: no content received before cancellation")
	}
}

func TestClaudeStreamError(t *testing.T) {
	srv := newClaudeStreamServer(t, []string{
		`event: error
data: {"type":"error","error":{"type":"rate_limit_error","message":"Rate limited"}}`,
	})
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)

	ch, err := cl.Stream(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}

	hadError := false
	for chunk := range ch {
		if chunk.Err != nil {
			hadError = true
			if !strings.Contains(chunk.Err.Error(), ErrRateLimited.Error()) {
				t.Errorf("unexpected error: %v", chunk.Err)
			}
		}
	}
	if !hadError {
		t.Error("expected error in stream")
	}
}

func TestClaudeDefaultConfig(t *testing.T) {
	cfg := DefaultClaudeConfig()
	if cfg.BaseURL != "https://api.anthropic.com" {
		t.Errorf("baseURL: got=%s", cfg.BaseURL)
	}
	if cfg.DefaultModel != "claude-sonnet-4-20250514" {
		t.Errorf("model: got=%s", cfg.DefaultModel)
	}
	if cfg.DefaultMaxTokens != 4096 {
		t.Errorf("maxTokens: got=%d", cfg.DefaultMaxTokens)
	}
}

func TestClaudeOptions(t *testing.T) {
	cl := NewClaude(
		WithClaudeAPIKey("key-123"),
		WithClaudeBaseURL("https://custom.example.com"),
		WithClaudeModel("custom-model"),
		WithClaudeTimeout(30*time.Second),
	)

	if cl.config.APIKey != "key-123" {
		t.Errorf("apiKey: got=%s", cl.config.APIKey)
	}
	if cl.config.BaseURL != "https://custom.example.com" {
		t.Errorf("baseURL: got=%s", cl.config.BaseURL)
	}
	if cl.config.DefaultModel != "custom-model" {
		t.Errorf("model: got=%s", cl.config.DefaultModel)
	}
}

func TestClaudeStreamingEndpointCheck(t *testing.T) {
	// Verify stream requests include ?stream=true and proper headers
	var method, path, apiKey string
	var streamField bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		apiKey = r.Header.Get("x-api-key")

		var req anthropicRequest
		json.NewDecoder(r.Body).Decode(&req)
		streamField = req.Stream

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cl := newTestClaude(t, srv.URL)
	ch, _ := cl.Stream(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	// Drain
	for range ch {
	}

	if method != http.MethodPost {
		t.Errorf("method: got=%s", method)
	}
	if path != "/v1/messages" {
		t.Errorf("path: got=%s", path)
	}
	if apiKey != "test-key" {
		t.Errorf("apiKey: got=%s", apiKey)
	}
	if !streamField {
		t.Error("stream field should be true")
	}
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

type claudeTestHandler struct {
	statusCode int
	response   string
}

func newClaudeTestServer(t *testing.T, h claudeTestHandler) *httptest.Server {
	t.Helper()
	code := h.statusCode
	if code == 0 {
		code = http.StatusOK
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, _ = w.Write([]byte(h.response))
	}))
}

func newClaudeStreamServer(t *testing.T, events []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("response writer does not support flushing")
		}

		for _, event := range events {
			if event == "" {
				// Empty line = SSE event separator
				_, _ = w.Write([]byte("\n"))
			} else {
				_, _ = w.Write([]byte(event + "\n"))
			}
			flusher.Flush()
		}
	}))
}

func newTestClaude(t *testing.T, baseURL string) *Claude {
	t.Helper()
	return NewClaude(
		WithClaudeBaseURL(baseURL),
		WithClaudeAPIKey("test-key"),
	)
}
