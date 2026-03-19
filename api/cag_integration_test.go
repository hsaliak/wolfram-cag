package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"wolfapi/client"
	"wolfapi/config"
)

func TestServiceEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		method     string
		invoke     func(*Service) (GenericResponse, []byte, error)
		bodyChecks []string
		queryCheck string
	}{
		{
			name:   "context",
			path:   "/api/cag/v1/WolframAlphaContext",
			method: http.MethodPost,
			invoke: func(s *Service) (GenericResponse, []byte, error) {
				return s.Context(context.Background(), ContextRequest{Context: "hello"})
			},
			bodyChecks: []string{`"context":"hello"`},
		},
		{
			name:   "hints",
			path:   "/api/cag/v1/WolframLanguageHints",
			method: http.MethodPost,
			invoke: func(s *Service) (GenericResponse, []byte, error) {
				return s.Hints(context.Background(), HintsRequest{Context: "hint me"})
			},
			bodyChecks: []string{`"context":"hint me"`},
		},
		{
			name:   "compute-with-opts",
			path:   "/api/cag/v1/WolframLanguageCompute",
			method: http.MethodPost,
			invoke: func(s *Service) (GenericResponse, []byte, error) {
				return s.Compute(context.Background(), ComputeRequest{Code: "2+2"}, ComputeOptions{TimeConstraint: 5, Line: 1, MaxChars: 100})
			},
			bodyChecks: []string{`"code":"2+2"`, `"timeConstraint":5`, `"line":1`, `"maxChars":100`},
		},
		{
			name:   "result-with-opts",
			path:   "/api/cag/v1/WolframAlphaResult",
			method: http.MethodGet,
			invoke: func(s *Service) (GenericResponse, []byte, error) {
				return s.Result(context.Background(), "integrate x", ResultOptions{Units: "metric", Format: "plaintext"})
			},
			queryCheck: "format=plaintext&input=integrate+x&units=metric",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path {
					t.Fatalf("path mismatch: got %s want %s", r.URL.Path, tc.path)
				}
				if r.Method != tc.method {
					t.Fatalf("method mismatch: got %s want %s", r.Method, tc.method)
				}
				if got := r.Header.Get("Authorization"); got != "test-key" {
					t.Fatalf("auth header mismatch: got %q", got)
				}

				if tc.method == http.MethodPost {
					if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
						t.Fatalf("content-type mismatch: got %q", got)
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("read body: %v", err)
					}
					for _, check := range tc.bodyChecks {
						if !strings.Contains(string(body), check) {
							t.Fatalf("body missing %q in %s", check, string(body))
						}
					}
				}

				if tc.method == http.MethodGet && tc.queryCheck != "" {
					if got := r.URL.RawQuery; got != tc.queryCheck {
						t.Fatalf("query mismatch: got %q want %q", got, tc.queryCheck)
					}
				}

				_ = json.NewEncoder(w).Encode(GenericResponse{Result: "ok"})
			}))
			defer server.Close()

			cfg := config.Config{
				APIKey:      "test-key",
				BaseURL:     server.URL + "/api/cag/v1",
				Output:      "json",
				TimeoutSecs: 2,
			}
			svc := New(client.New(cfg))

			resp, _, err := tc.invoke(svc)
			if err != nil {
				t.Fatalf("invoke returned error: %v", err)
			}
			if resp.Result != "ok" {
				t.Fatalf("unexpected response result: %q", resp.Result)
			}
		})
	}
}

func TestResultPlaintextFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cag/v1/WolframAlphaResult" {
			t.Fatalf("path mismatch: got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("format"); got != "plaintext" {
			t.Fatalf("format query mismatch: got %q", got)
		}
		_, _ = w.Write([]byte("Query: weather in Boston\nResult: 18 °C"))
	}))
	defer server.Close()

	cfg := config.Config{
		APIKey:      "test-key",
		BaseURL:     server.URL + "/api/cag/v1",
		Output:      "text",
		TimeoutSecs: 2,
	}
	svc := New(client.New(cfg))

	resp, raw, err := svc.Result(context.Background(), "weather in boston", ResultOptions{Format: "plaintext"})
	if err != nil {
		t.Fatalf("expected plaintext fallback, got error: %v", err)
	}
	if !strings.Contains(resp.Result, "Query: weather in Boston") {
		t.Fatalf("unexpected fallback result: %q", resp.Result)
	}
	if !strings.Contains(string(raw), "Result: 18") {
		t.Fatalf("unexpected raw body: %q", string(raw))
	}
}
