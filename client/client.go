package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"wolfapi/config"
	"wolfapi/errs"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	verbose    bool
}

func New(cfg config.Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: time.Duration(cfg.TimeoutSecs) * time.Second},
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:     cfg.APIKey,
		verbose:    cfg.Verbose,
	}
}

func (c *Client) Do(ctx context.Context, method, endpointPath string, query url.Values, payload any) ([]byte, error) {
	u := c.baseURL + "/" + strings.TrimLeft(endpointPath, "/")
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, errs.EncodeError{Err: err}
		}
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, errs.InvalidArgsError{Msg: err.Error()}
	}

	req.Header.Set("Authorization", c.apiKey)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.MapRequestError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.NetworkError{Err: err}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errs.HTTPStatusError{Code: resp.StatusCode, Body: strings.TrimSpace(string(body))}
	}

	return body, nil
}

func DecodeJSON(data []byte, out any) error {
	if err := json.Unmarshal(data, out); err != nil {
		return errs.DecodeError{Err: err}
	}
	return nil
}

func (c *Client) Verbose() bool {
	return c.verbose
}
