package config

import (
	"os"
	"strings"

	"wolfapi/errs"
)

const (
	DefaultBaseURL = "https://services.wolfram.com/api/cag/v1"
	DefaultOutput  = "text"
	DefaultTimeout = 30
	EnvAPIKey      = "WOLFRAM_APP_ID"
)

type Config struct {
	APIKey      string
	BaseURL     string
	Output      string
	TimeoutSecs int
	Verbose     bool
	Workers     int
}

type Inputs struct {
	APIKey      string
	BaseURL     string
	Output      string
	TimeoutSecs int
	Verbose     bool
	Workers     int
}

func Resolve(in Inputs) (Config, error) {
	apiKey := strings.TrimSpace(in.APIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv(EnvAPIKey))
	}
	if apiKey == "" {
		return Config{}, errs.ErrMissingAPIKey
	}

	baseURL := strings.TrimSpace(in.BaseURL)
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	output := strings.ToLower(strings.TrimSpace(in.Output))
	if output == "" {
		output = DefaultOutput
	}

	timeout := in.TimeoutSecs
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	workers := in.Workers
	if workers <= 0 {
		workers = 1
	}

	return Config{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Output:      output,
		TimeoutSecs: timeout,
		Verbose:     in.Verbose,
		Workers:     workers,
	}, nil
}
