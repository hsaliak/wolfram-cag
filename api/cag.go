package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"wolfapi/client"
)

type ContextRequest struct {
	Context string `json:"context"`
}

type HintsRequest struct {
	Context string `json:"context"`
}

type ComputeRequest struct {
	Code string `json:"code"`
}

type ComputeOptions struct {
	TimeConstraint int `json:"timeConstraint,omitempty"`
	Line           int `json:"line,omitempty"`
	MaxChars       int `json:"maxChars,omitempty"`
}

type ResultOptions struct {
	Assumption string
	Format     string
	Units      string
	Location   string
	LatLong    string
	Timeout    string
	MaxWidth   string
}

type GenericResponse struct {
	Result  string `json:"result"`
	Code    int    `json:"code,omitempty"`
	Success *bool  `json:"success,omitempty"`
	UUID    string `json:"uuid,omitempty"`
}

type ComputePayload struct {
	Code string `json:"code"`
	ComputeOptions
}

type Service struct {
	client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{client: c}
}

func (s *Service) Context(ctx context.Context, req ContextRequest) (GenericResponse, []byte, error) {
	body, err := s.client.Do(ctx, http.MethodPost, "/WolframAlphaContext", nil, req)
	if err != nil {
		return GenericResponse{}, nil, err
	}

	var resp GenericResponse
	if err := client.DecodeJSON(body, &resp); err != nil {
		// WolframAlphaResult may return plaintext (for example with
		// format=plaintext). In that case, surface the raw body as Result.
		fallback := strings.TrimSpace(string(body))
		if fallback == "" {
			return GenericResponse{}, nil, err
		}
		resp = GenericResponse{Result: fallback}
	}

	return resp, body, nil
}

func (s *Service) Hints(ctx context.Context, req HintsRequest) (GenericResponse, []byte, error) {
	body, err := s.client.Do(ctx, http.MethodPost, "/WolframLanguageHints", nil, req)
	if err != nil {
		return GenericResponse{}, nil, err
	}

	var resp GenericResponse
	if err := client.DecodeJSON(body, &resp); err != nil {
		// WolframAlphaResult may return plaintext (for example with
		// format=plaintext). In that case, surface the raw body as Result.
		fallback := strings.TrimSpace(string(body))
		if fallback == "" {
			return GenericResponse{}, nil, err
		}
		resp = GenericResponse{Result: fallback}
	}

	return resp, body, nil
}

func (s *Service) Compute(ctx context.Context, req ComputeRequest, opts ComputeOptions) (GenericResponse, []byte, error) {
	body, err := s.client.Do(ctx, http.MethodPost, "/WolframLanguageCompute", nil, ComputePayload{Code: req.Code, ComputeOptions: opts})
	if err != nil {
		return GenericResponse{}, nil, err
	}

	var resp GenericResponse
	if err := client.DecodeJSON(body, &resp); err != nil {
		// WolframAlphaResult may return plaintext (for example with
		// format=plaintext). In that case, surface the raw body as Result.
		fallback := strings.TrimSpace(string(body))
		if fallback == "" {
			return GenericResponse{}, nil, err
		}

		resp = GenericResponse{Result: fallback}
	}

	return resp, body, nil
}

func (s *Service) Result(ctx context.Context, input string, opts ResultOptions) (GenericResponse, []byte, error) {
	query, err := BuildResultQuery(input, opts)
	if err != nil {
		return GenericResponse{}, nil, err
	}

	body, err := s.client.Do(ctx, http.MethodGet, "/WolframAlphaResult", query, nil)
	if err != nil {
		return GenericResponse{}, nil, err
	}

	var resp GenericResponse
	if err := client.DecodeJSON(body, &resp); err != nil {
		// WolframAlphaResult may return plaintext (for example with
		// format=plaintext). In that case, surface the raw body as Result.
		fallback := strings.TrimSpace(string(body))
		if fallback == "" {
			return GenericResponse{}, nil, err
		}

		resp = GenericResponse{Result: fallback}
	}

	return resp, body, nil
}

func BuildResultQuery(input string, opts ResultOptions) (url.Values, error) {
	trimmedInput := strings.TrimSpace(input)
	if trimmedInput == "" {
		return nil, fmt.Errorf("input is required")
	}

	query := url.Values{}
	query.Set("input", trimmedInput)

	setIfNotEmpty(query, "assumption", opts.Assumption)
	setIfNotEmpty(query, "format", opts.Format)
	setIfNotEmpty(query, "units", opts.Units)
	setIfNotEmpty(query, "location", opts.Location)
	setIfNotEmpty(query, "latlong", opts.LatLong)
	setIfNotEmpty(query, "timeout", opts.Timeout)
	setIfNotEmpty(query, "maxwidth", opts.MaxWidth)

	return query, nil
}

func setIfNotEmpty(v url.Values, k, val string) {
	if val == "" {
		return
	}
	v.Set(k, val)
}
