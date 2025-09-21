package transform

import (
	"bytes"
	"testing"

	"golang.org/x/text/transform"
)

func TestNewlineTransformer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Windows line endings",
			input:    "line1\r\nline2\r\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "Mixed line endings",
			input:    "line1\r\nline2\nline3\r\n",
			expected: "line1\nline2\nline3\n",
		},
		{
			name:     "Unix line endings only",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "No line endings",
			input:    "single line",
			expected: "single line",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only carriage returns",
			input:    "line1\rline2\r",
			expected: "line1\rline2\r",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := WindowsToUnixNewlines()

			// Transform the input
			result, _, err := transform.String(transformer, tt.input)
			if err != nil {
				t.Fatalf("Transform failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNewlineTransformerWithReader(t *testing.T) {
	input := "Windows\r\nline\r\nendings\r\n"
	expected := "Windows\nline\nendings\n"

	transformer := WindowsToUnixNewlines()
	reader := transform.NewReader(bytes.NewReader([]byte(input)), transformer)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		t.Fatalf("Reading from transformer failed: %v", err)
	}

	result := buf.String()
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNewlineTransformerWithWriter(t *testing.T) {
	input := "Windows\r\nline\r\nendings\r\n"
	expected := "Windows\nline\nendings\n"

	transformer := WindowsToUnixNewlines()
	var buf bytes.Buffer
	writer := transform.NewWriter(&buf, transformer)

	_, err := writer.Write([]byte(input))
	if err != nil {
		t.Fatalf("Writing to transformer failed: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Closing transformer writer failed: %v", err)
	}

	result := buf.String()
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
