package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// FraudClient defines an interface for calling external fraud detection.
type FraudClient interface {
	Evaluate(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error)
}

// HTTPFraudClient is a simple HTTP implementation. It can be swapped with a mock in tests.
type HTTPFraudClient struct {
	url string
	c   *http.Client
}

func NewHTTPFraudClient(url string) FraudClient {
	return &HTTPFraudClient{url: url, c: &http.Client{Timeout: 5 * time.Second}}
}

func (h *HTTPFraudClient) Evaluate(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if h.url == "" {
		// no external API configured, return a benign response
		return map[string]interface{}{"score": 0.1, "recommendation": "approve"}, nil
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", h.url, bytesReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// helper to avoid importing bytes in multiple files
func bytesReader(b []byte) *bytesReaderT { return &bytesReaderT{b, 0} }

// a tiny implementation to avoid extra imports — used only inside this example
type bytesReaderT struct {
	b []byte
	i int
}

func (r *bytesReaderT) Read(p []byte) (n int, err error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n = copy(p, r.b[r.i:])
	r.i += n
	return
}
func (r *bytesReaderT) Close() error { return nil }
