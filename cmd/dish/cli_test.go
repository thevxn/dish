package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printHelp()

	if err := w.Close(); err != nil {
		t.Errorf("pipe close: %v", err)
	}

	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, "Usage: dish [FLAGS] SOURCE") {
		t.Errorf("help output missing usage line")
	}
	if !strings.Contains(output, "A lightweight, one-shot socket checker") {
		t.Errorf("help output missing description")
	}
	if !strings.Contains(output, "SOURCE must be a file path") {
		t.Errorf("help output missing source description")
	}
	if !strings.Contains(output, "Use the `-h` flag") {
		t.Errorf("help output missing -h flag info")
	}
}
