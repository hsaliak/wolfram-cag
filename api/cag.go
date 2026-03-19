package api

import (
	"context"
	"net/http"
	"net/url"

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

type GenericResponse struct {
	Result  string `json:"result"`
	Code    int    `json:"code,omitempty"`
	Success *bool  `json:"success,omitempty"`
	UUID    string `json:"uuid,omitempty"`
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
		return GenericResponse{}, nil, err
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
		return GenericResponse{}, nil, err
	}

	return resp, body, nil
}

func (s *Service) Compute(ctx context.Context, req ComputeRequest) (GenericResponse, []byte, error) {
	body, err := s.client.Do(ctx, http.MethodPost, "/WolframLanguageCompute", nil, req)
	if err != nil {
		return GenericResponse{}, nil, err
	}

	var resp GenericResponse
	if err := client.DecodeJSON(body, &resp); err != nil {
		return GenericResponse{}, nil, err
	}

	return resp, body, nil
}

func (s *Service) Result(ctx context.Context, input string) (GenericResponse, []byte, error) {
	query := url.Values{}
	query.Set("input", input)

	body, err := s.client.Do(ctx, http.MethodGet, "/WolframAlphaResult", query, nil)
	if err != nil {
		return GenericResponse{}, nil, err
	}

	var resp GenericResponse
	if err := client.DecodeJSON(body, &resp); err != nil {
		return GenericResponse{}, nil, err
	}

	return resp, body, nil
}
