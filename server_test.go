package main

import (
	"strings"
	"testing"
)

// TestAsciiArtValidBanner tests that AsciiArt returns non-empty output for valid inputs
func TestAsciiArtValidBanner(t *testing.T) {
	result := AsciiArt("Hello", "standard")
	if result == "" {
		t.Error("Expected non-empty result for valid input and banner, got empty string")
	}
}

// TestAsciiArtShadowBanner tests the shadow banner
func TestAsciiArtShadowBanner(t *testing.T) {
	result := AsciiArt("Hello", "shadow")
	if result == "" {
		t.Error("Expected non-empty result for shadow banner, got empty string")
	}
}

// TestAsciiArtThinkertoyBanner tests the thinkertoy banner
func TestAsciiArtThinkertoyBanner(t *testing.T) {
	result := AsciiArt("Hello", "thinkertoy")
	if result == "" {
		t.Error("Expected non-empty result for thinkertoy banner, got empty string")
	}
}

// TestAsciiArtInvalidBanner tests that AsciiArt returns empty string for a non-existent banner
func TestAsciiArtInvalidBanner(t *testing.T) {
	result := AsciiArt("Hello", "nonexistent")
	if result != "" {
		t.Error("Expected empty result for invalid banner, got non-empty string")
	}
}

// TestAsciiArtEmptyInput tests that empty input returns empty output
func TestAsciiArtEmptyInput(t *testing.T) {
	result := AsciiArt("", "standard")
	if result != "" {
		t.Errorf("Expected empty result for empty input, got: %q", result)
	}
}

// TestAsciiArtNewline tests that \n in input produces multi-line output
func TestAsciiArtNewline(t *testing.T) {
	result := AsciiArt("Hi\\nHi", "standard")
	if result == "" {
		t.Error("Expected non-empty result for multi-line input")
	}
	// Should contain more lines than a single word
	lines := strings.Split(result, "\n")
	if len(lines) <= 8 {
		t.Errorf("Expected more than 8 lines for multi-line input, got %d", len(lines))
	}
}

// TestAsciiArtOutputHasEightLinesPerWord tests that each word produces exactly 8 lines of art
func TestAsciiArtOutputHasEightLinesPerWord(t *testing.T) {
	result := AsciiArt("Hi", "standard")
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("Expected 8 lines of output for single word, got %d", len(lines))
	}
}
