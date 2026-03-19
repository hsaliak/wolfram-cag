package config

import "testing"

func TestResolve_UsesCLIKeyOverEnv(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-key")

	cfg, err := Resolve(Inputs{APIKey: "cli-key"})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if cfg.APIKey != "cli-key" {
		t.Fatalf("expected cli key, got %q", cfg.APIKey)
	}
}

func TestResolve_UsesEnvKeyWhenCLIEmpty(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-key")

	cfg, err := Resolve(Inputs{})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if cfg.APIKey != "env-key" {
		t.Fatalf("expected env key, got %q", cfg.APIKey)
	}
}

func TestResolve_MissingKey(t *testing.T) {
	t.Setenv(EnvAPIKey, "")

	_, err := Resolve(Inputs{})
	if err == nil {
		t.Fatalf("expected error for missing key")
	}
}

func TestResolve_Defaults(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-key")

	cfg, err := Resolve(Inputs{})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if cfg.BaseURL != DefaultBaseURL {
		t.Fatalf("base url default mismatch: %q", cfg.BaseURL)
	}
	if cfg.Output != DefaultOutput {
		t.Fatalf("output default mismatch: %q", cfg.Output)
	}
	if cfg.TimeoutSecs != DefaultTimeout {
		t.Fatalf("timeout default mismatch: %d", cfg.TimeoutSecs)
	}
	if cfg.Workers != 4 {
		t.Fatalf("workers default mismatch: %d", cfg.Workers)
	}
}

func TestResolve_NormalizesOutput(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-key")

	cfg, err := Resolve(Inputs{Output: " JSON "})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if cfg.Output != "json" {
		t.Fatalf("expected normalized output json, got %q", cfg.Output)
	}
}
