package config

import (
	"flag"
	"testing"
)

func TestNewConfig_DefaultsAndSource(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{"source.json"}

	cfg, err := NewConfig(fs, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Source != "source.json" {
		t.Errorf("expected Source to be 'source.json', got %q", cfg.Source)
	}
	if cfg.InstanceName != defaultInstanceName {
		t.Errorf("expected InstanceName %q, got %q", defaultInstanceName, cfg.InstanceName)
	}
	if cfg.ApiCacheDirectory != defaultApiCacheDir {
		t.Errorf("expected ApiCacheDirectory %q, got %q", defaultApiCacheDir, cfg.ApiCacheDirectory)
	}
	if cfg.Verbose != defaultVerbose {
		t.Errorf("expected Verbose %v, got %v", defaultVerbose, cfg.Verbose)
	}
}

func TestNewConfig_FlagsOverrideDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{
		"-name", "custom-dish",
		"-timeout", "42",
		"-verbose",
		"-hname", "X-Auth",
		"-hvalue", "secret",
		"-cache",
		"-cacheDir", "/tmp/cache",
		"-cacheTTL", "99",
		"-target", "http://push",
		"-telegramBotToken", "token",
		"-telegramChatID", "chatid",
		"-updateURL", "http://api",
		"-webhookURL", "http://webhook",
		"-textNotifySuccess",
		"-machineNotifySuccess",
		"mysource.json",
	}

	cfg, err := NewConfig(fs, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.InstanceName != "custom-dish" {
		t.Errorf("expected InstanceName 'custom-dish', got %q", cfg.InstanceName)
	}
	if cfg.TimeoutSeconds != 42 {
		t.Errorf("expected TimeoutSeconds 42, got %d", cfg.TimeoutSeconds)
	}
	if !cfg.Verbose {
		t.Errorf("expected Verbose true, got false")
	}
	if cfg.ApiHeaderName != "X-Auth" {
		t.Errorf("expected ApiHeaderName 'X-Auth', got %q", cfg.ApiHeaderName)
	}
	if cfg.ApiHeaderValue != "secret" {
		t.Errorf("expected ApiHeaderValue 'secret', got %q", cfg.ApiHeaderValue)
	}
	if !cfg.ApiCacheSockets {
		t.Errorf("expected ApiCacheSockets true, got false")
	}
	if cfg.ApiCacheDirectory != "/tmp/cache" {
		t.Errorf("expected ApiCacheDirectory '/tmp/cache', got %q", cfg.ApiCacheDirectory)
	}
	if cfg.ApiCacheTTLMinutes != 99 {
		t.Errorf("expected ApiCacheTTLMinutes 99, got %d", cfg.ApiCacheTTLMinutes)
	}
	if cfg.PushgatewayURL != "http://push" {
		t.Errorf("expected PushgatewayURL 'http://push', got %q", cfg.PushgatewayURL)
	}
	if cfg.TelegramBotToken != "token" {
		t.Errorf("expected TelegramBotToken 'token', got %q", cfg.TelegramBotToken)
	}
	if cfg.TelegramChatID != "chatid" {
		t.Errorf("expected TelegramChatID 'chatid', got %q", cfg.TelegramChatID)
	}
	if cfg.ApiURL != "http://api" {
		t.Errorf("expected ApiURL 'http://api', got %q", cfg.ApiURL)
	}
	if cfg.WebhookURL != "http://webhook" {
		t.Errorf("expected WebhookURL 'http://webhook', got %q", cfg.WebhookURL)
	}
	if !cfg.TextNotifySuccess {
		t.Errorf("expected TextNotifySuccess true, got false")
	}
	if !cfg.MachineNotifySuccess {
		t.Errorf("expected MachineNotifySuccess true, got false")
	}
	if cfg.Source != "mysource.json" {
		t.Errorf("expected Source 'mysource.json', got %q", cfg.Source)
	}
}

func TestNewConfig_NoSourceProvided(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{}

	_, err := NewConfig(fs, args)
	if err == nil {
		t.Fatal("expected error for no source provided, got nil")
	}
	if err != ErrNoSourceProvided {
		t.Fatalf("expected ErrNoSourceProvided, got %v", err)
	}
}

func TestNewConfig_InvalidFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{"-notaflag", "source.json"}

	_, err := NewConfig(fs, args)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
