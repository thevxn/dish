package config

import (
	"flag"
	"reflect"
	"testing"
)

func TestNewConfig_DefaultsAndSource(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{"source.json"}

	expected := &Config{
		InstanceName:         defaultInstanceName,
		ApiHeaderName:        defaultApiHeaderName,
		ApiHeaderValue:       defaultApiHeaderValue,
		ApiCacheSockets:      defaultApiCacheSockets,
		ApiCacheDirectory:    defaultApiCacheDir,
		ApiCacheTTLMinutes:   defaultApiCacheTTLMinutes,
		Source:               "source.json",
		Verbose:              defaultVerbose,
		PushgatewayURL:       defaultPushgatewayURL,
		TelegramBotToken:     defaultTelegramBotToken,
		TelegramChatID:       defaultTelegramChatID,
		TimeoutSeconds:       defaultTimeoutSeconds,
		ApiURL:               defaultApiURL,
		WebhookURL:           defaultWebhookURL,
		TextNotifySuccess:    defaultTextNotifySuccess,
		MachineNotifySuccess: defaultMachineNotifySuccess,
	}

	if blank, err := NewConfig(nil, []string{}); err == nil || blank != nil {
		t.Fatalf("unexpected behaviour, err should not be nil, output should be nil")
	}

	actual, err := NewConfig(fs, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
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

	expected := &Config{
		InstanceName:         "custom-dish",
		TimeoutSeconds:       42,
		Verbose:              true,
		ApiHeaderName:        "X-Auth",
		ApiHeaderValue:       "secret",
		ApiCacheSockets:      true,
		ApiCacheDirectory:    "/tmp/cache",
		ApiCacheTTLMinutes:   99,
		PushgatewayURL:       "http://push",
		TelegramBotToken:     "token",
		TelegramChatID:       "chatid",
		ApiURL:               "http://api",
		WebhookURL:           "http://webhook",
		TextNotifySuccess:    true,
		MachineNotifySuccess: true,
		Source:               "mysource.json",
	}

	actual, err := NewConfig(fs, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
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
