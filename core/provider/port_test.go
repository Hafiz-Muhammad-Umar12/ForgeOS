package provider

import (
	"context"
	"testing"
)

func TestRoles(t *testing.T) {
	if RoleUser != "user" {
		t.Errorf("RoleUser=%s", RoleUser)
	}
	if RoleAssistant != "assistant" {
		t.Errorf("RoleAssistant=%s", RoleAssistant)
	}
	if RoleSystem != "system" {
		t.Errorf("RoleSystem=%s", RoleSystem)
	}
}

func TestCompletionRequestDefaults(t *testing.T) {
	req := CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "hello"}},
	}
	if req.Model != "" {
		t.Errorf("model should be empty")
	}
	if len(req.Messages) != 1 {
		t.Errorf("messages=%d", len(req.Messages))
	}
}

func TestCompletionResponseFields(t *testing.T) {
	resp := CompletionResponse{
		Message:      Message{Role: RoleAssistant, Content: "Hi!"},
		FinishReason: "end_turn",
		Usage:        Usage{InputTokens: 10, OutputTokens: 5},
	}
	if resp.Message.Role != RoleAssistant {
		t.Errorf("role=%s", resp.Message.Role)
	}
	if resp.Message.Content != "Hi!" {
		t.Errorf("content=%s", resp.Message.Content)
	}
	if resp.FinishReason != "end_turn" {
		t.Errorf("reason=%s", resp.FinishReason)
	}
	if resp.Usage.InputTokens != 10 {
		t.Errorf("input=%d", resp.Usage.InputTokens)
	}
}

func TestStreamChunk(t *testing.T) {
	t.Run("content chunk", func(t *testing.T) {
		chunk := StreamChunk{Content: "Hello"}
		if chunk.Content != "Hello" {
			t.Errorf("content=%s", chunk.Content)
		}
		if chunk.Done {
			t.Error("done should be false")
		}
	})

	t.Run("done chunk with usage", func(t *testing.T) {
		chunk := StreamChunk{Done: true, Usage: Usage{InputTokens: 5, OutputTokens: 10}}
		if !chunk.Done {
			t.Error("done should be true")
		}
		if chunk.Usage.OutputTokens != 10 {
			t.Errorf("output=%d", chunk.Usage.OutputTokens)
		}
	})

	t.Run("error chunk", func(t *testing.T) {
		chunk := StreamChunk{Err: ErrRateLimited}
		if chunk.Err != ErrRateLimited {
			t.Errorf("err=%v", chunk.Err)
		}
	})
}

func TestCapabilitiesDefaults(t *testing.T) {
	caps := Capabilities{
		Provider:  "test",
		Models:    []string{"m1", "m2"},
		Streaming: true,
		MaxTokens: 100,
	}
	if caps.Provider != "test" {
		t.Errorf("provider=%s", caps.Provider)
	}
	if len(caps.Models) != 2 {
		t.Errorf("models=%v", caps.Models)
	}
	if !caps.Streaming {
		t.Error("streaming should be true")
	}
	if caps.MaxTokens != 100 {
		t.Errorf("maxtokens=%d", caps.MaxTokens)
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		err   error
		label string
	}{
		{ErrRateLimited, "ErrRateLimited"},
		{ErrContextWindowExceeded, "ErrContextWindowExceeded"},
		{ErrProviderUnavailable, "ErrProviderUnavailable"},
		{ErrAuthFailed, "ErrAuthFailed"},
		{ErrBadRequest, "ErrBadRequest"},
		{ErrStreamInterrupted, "ErrStreamInterrupted"},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("error is nil")
			}
		})
	}
}

func TestFakeProviderComplete(t *testing.T) {
	fp := NewFakeProvider(CompletionResponse{
		Message:      Message{Role: RoleAssistant, Content: "Hello!"},
		FinishReason: "end_turn",
		Usage:        Usage{InputTokens: 5, OutputTokens: 3},
	})

	resp, err := fp.Complete(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if resp.Message.Content != "Hello!" {
		t.Errorf("content: got=%s", resp.Message.Content)
	}
	if fp.RequestCount.Load() != 1 {
		t.Errorf("count=%d", fp.RequestCount.Load())
	}
	if len(fp.RequestsReceived) != 1 {
		t.Errorf("requests=%d", len(fp.RequestsReceived))
	}
}

func TestFakeProviderCompleteExhausted(t *testing.T) {
	fp := NewFakeProvider()
	// Consume the default response
	_, _ = fp.Complete(context.Background(), CompletionRequest{})
	// Second call should fail
	_, err := fp.Complete(context.Background(), CompletionRequest{})
	if err == nil {
		t.Fatal("expected error on exhausted responses")
	}
}

func TestFakeProviderStream(t *testing.T) {
	fp := NewFakeProvider(CompletionResponse{
		Message:      Message{Role: RoleAssistant, Content: "Hello world!"},
		FinishReason: "end_turn",
		Usage:        Usage{InputTokens: 5, OutputTokens: 3},
	})

	ch, err := fp.Stream(context.Background(), CompletionRequest{
		Messages: []Message{{Role: RoleUser, Content: "Hi"}},
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
		if chunk.Done {
			if chunk.Usage.OutputTokens != 3 {
				t.Errorf("usage: got=%d", chunk.Usage.OutputTokens)
			}
		}
	}
	if full != "Hello world!" {
		t.Errorf("streamed content: got=%s", full)
	}
}

func TestFakeProviderCapabilities(t *testing.T) {
	fp := NewFakeProvider()
	caps := fp.Capabilities()
	if caps.Provider != "fake" {
		t.Errorf("provider=%s", caps.Provider)
	}
}

func TestFakeProviderCustomCompleteFunc(t *testing.T) {
	fp := NewFakeProvider()
	fp.CompleteFunc = func(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
		return CompletionResponse{
			Message:      Message{Role: RoleAssistant, Content: "custom"},
			FinishReason: "end_turn",
		}, nil
	}

	resp, err := fp.Complete(context.Background(), CompletionRequest{})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if resp.Message.Content != "custom" {
		t.Errorf("content: got=%s", resp.Message.Content)
	}
}

func TestFakeProviderRequestRecording(t *testing.T) {
	fp := NewFakeProvider()

	req := CompletionRequest{
		Model:       "test-model",
		Messages:    []Message{{Role: RoleUser, Content: "record me"}},
		System:      "be helpful",
		MaxTokens:   100,
		Temperature: 0.5,
	}

	_, _ = fp.Complete(context.Background(), req)

	if len(fp.RequestsReceived) != 1 {
		t.Fatalf("expected 1 request, got %d", len(fp.RequestsReceived))
	}
	recorded := fp.RequestsReceived[0]
	if recorded.Model != "test-model" {
		t.Errorf("model: got=%s", recorded.Model)
	}
	if len(recorded.Messages) != 1 || recorded.Messages[0].Content != "record me" {
		t.Errorf("messages: got=%v", recorded.Messages)
	}
	if recorded.System != "be helpful" {
		t.Errorf("system: got=%s", recorded.System)
	}
	if recorded.MaxTokens != 100 {
		t.Errorf("maxTokens: got=%d", recorded.MaxTokens)
	}
	if recorded.Temperature != 0.5 {
		t.Errorf("temperature: got=%f", recorded.Temperature)
	}
}

func TestFakeProviderCustomCapabilities(t *testing.T) {
	fp := NewFakeProvider()
	custom := Capabilities{
		Provider:  "custom",
		Models:    []string{"custom-model"},
		Streaming: false,
		MaxTokens: 2048,
	}
	fp.ConfiguredCapabilities = &custom

	caps := fp.Capabilities()
	if caps.Provider != "custom" {
		t.Errorf("provider: got=%s", caps.Provider)
	}
	if caps.Streaming {
		t.Error("streaming should be false")
	}
}

func TestSplitContent(t *testing.T) {
	tests := []struct {
		input     string
		chunkSize int
		expected  int
	}{
		{"hello", 2, 3},
		{"hi", 10, 1},
		{"", 5, 1},
	}
	for _, tt := range tests {
		chunks := splitContent(tt.input, tt.chunkSize)
		var joined string
		for _, c := range chunks {
			joined += c
		}
		if joined != tt.input {
			t.Errorf("split/join: got=%s want=%s", joined, tt.input)
		}
	}
}
