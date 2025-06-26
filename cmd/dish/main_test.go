package main

import (
	"bytes"
	"flag"
	"os"
	"testing"
)

func TestRun_InvalidFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	stderr := &bytes.Buffer{}

	code := run(fs, []string{"-notaflag"}, os.Stdout, stderr)
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

func TestRun_NoArgs(t *testing.T) {
	fs := flag.NewFlagSet("no_args", flag.ContinueOnError)
	stderr := &bytes.Buffer{}

	code := run(fs, []string{}, os.Stdout, stderr)
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRun_ValidSockets(t *testing.T) {
	fs := flag.NewFlagSet("valid_sockets", flag.ContinueOnError)
	stderr := &bytes.Buffer{}
	tmpfile := testFile(t, "test_sockets.json", []byte(testSocketsValid))

	code := run(fs, []string{tmpfile}, os.Stdout, stderr)
	if code != 0 {
		t.Errorf("expected exit code 0 got %d", code)
	}
}

func TestRun_InvalidSockets(t *testing.T) {
	fs := flag.NewFlagSet("invalid_sockets", flag.ContinueOnError)
	stderr := &bytes.Buffer{}
	tmpfile := testFile(t, "test_sockets.json", []byte(testSocketsSomeInvalid))

	code := run(fs, []string{tmpfile}, os.Stdout, stderr)
	if code != 4 {
		t.Errorf("expected exit code 4 got %d", code)
	}
}

func TestRun_InvalidSource(t *testing.T) {
	fs := flag.NewFlagSet("invalid_source", flag.ContinueOnError)
	stderr := &bytes.Buffer{}

	code := run(fs, []string{""}, os.Stdout, stderr)
	if code != 3 {
		t.Errorf("expected exit code 3 got %d", code)
	}
}
