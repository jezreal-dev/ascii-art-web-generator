package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAsciiArtValidBanner tests that AsciiArt returns non-empty output for valid inputs
func TestAsciiArtValidBanner(t *testing.T) {
	result, err := AsciiArt("Hello", "standard")
	if err != nil {
    	t.Errorf("Expected no error, got: %v", err)
	}
	if result == "" {
    	t.Error("Expected non-empty result, got empty string")
	}
}

// TestAsciiArtShadowBanner tests the shadow banner
func TestAsciiArtShadowBanner(t *testing.T) {
	result, err := AsciiArt("Hello", "shadow")
	if err != nil {
    	t.Errorf("Expected no error, got: %v", err)
	}
	if result == "" {
    	t.Error("Expected non-empty result, got empty string")
	}
}

// TestAsciiArtThinkertoyBanner tests the thinkertoy banner
func TestAsciiArtThinkertoyBanner(t *testing.T) {
	result, err := AsciiArt("Hello", "thinkertoy")
	if err != nil {
    	t.Errorf("Expected no error, got: %v", err)
	}
	if result == "" {
    	t.Error("Expected non-empty result, got empty string")
	}
}

// TestAsciiArtInvalidBanner tests that AsciiArt returns empty string for a non-existent banner
func TestAsciiArtInvalidBanner(t *testing.T) {
	result, err := AsciiArt("Hello", "nonexistent")
	if err == nil {
    	t.Error("Expected error for invalid banner, got nil")
	}
	if result != "" {
    	t.Error("Expected empty result for invalid banner")
	}
}

// TestAsciiArtEmptyInput tests that empty input returns empty output
func TestAsciiArtEmptyInput(t *testing.T) {
	result, err := AsciiArt("", "standard")
	if err != nil {
    	t.Error("Expected no error for empty input")
	}
	if result != "" {
    	t.Error("Expected empty result for empty input")
	}
}

// TestAsciiArtNewline tests that \n in input produces multi-line output
func TestAsciiArtNewline(t *testing.T) {
	result, err := AsciiArt("Hi\\nHi", "standard")
	if err != nil {
    	t.Errorf("Expected no error for valid input, got: %v", err)
	}
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
	result, err := AsciiArt("Hi", "standard")
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("Expected 8 lines of output for single word, got %d", len(lines))
	}
	if err != nil {
    	t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

func TestAsciiArtSpecialCharacters(t *testing.T) {
	result, err := AsciiArt("123??", "standard")
	if err != nil {
    	t.Errorf("Expected no error, got: %v", err)
	}
	if result == "" {
    	t.Error("Expected non-empty result, got empty string")
	}
}

// 
func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	homeHandler(rr, req)
	if rr.Code != http.StatusOK {
    	t.Errorf("Expected 200, got %d", rr.Code)
	}
}

func TestHomeHandlerInvalidPath(t *testing.T) {
	req := httptest.NewRequest("GET", "/invalid", nil)
	rr := httptest.NewRecorder()
	homeHandler(rr, req)
	if rr.Code != http.StatusNotFound {
    	t.Errorf("Expected 404, got %d", rr.Code)
	}
}

func TestAsciiArtHandlerReturns200(t *testing.T) {
	body := strings.NewReader("text=Hello&banner=standard")
	req := httptest.NewRequest("POST", "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	asciiArtHandler(rr, req)
	if rr.Code != http.StatusOK {
    	t.Errorf("Expected 200, got %d", rr.Code)
	}
}


func TestAsciiArtHandlerEmptyText(t *testing.T) {
	body := strings.NewReader("text=&banner=standard")
	req := httptest.NewRequest("POST", "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	asciiArtHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
    	t.Errorf("Expected 400, got %d", rr.Code)
	}
}


func TestAsciiArtHandlerInvalidBanner(t *testing.T) {
	body := strings.NewReader("text=Hello&banner=fakebanner")
	req := httptest.NewRequest("POST", "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	asciiArtHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
    	t.Errorf("Expected 400, got %d", rr.Code)
	}
}

func TestAsciiArtHandlerWrongMethod(t *testing.T) {
	body := strings.NewReader("text=Hello&banner=fakebanner")
	req := httptest.NewRequest("GET", "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	asciiArtHandler(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
    	t.Errorf("Expected 405, got %d", rr.Code)
	}
}

func TestAsciiArtHandlerMultiLine(t *testing.T) {
	body := strings.NewReader("text=Hello%5CnWorld&banner=standard")
	req := httptest.NewRequest("POST", "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	asciiArtHandler(rr, req)
	if rr.Code != http.StatusOK {
    	t.Errorf("Expected 200, got %d", rr.Code)
	}
}

func TestAsciiArtInvalidCharacter(t *testing.T) {
	_, err := AsciiArt("Hello 😊", "standard")
	if err == nil {
		t.Error("Expected error for non-ASCII emoji character, got nil")
	}
	if err != nil && err.Error() != "Invalid character in input" {
		t.Errorf("Expected 'Invalid character in input' error, got: %v", err)
	}
}

func TestAsciiArtActualNewline(t *testing.T) {
	result, err := AsciiArt("Hello\nWorld", "standard")
	if err != nil {
		t.Errorf("Expected no error for actual newline, got: %v", err)
	}
	if result == "" {
		t.Error("Expected non-empty result for actual newline")
	}
}

func TestAsciiArtOnlyNewlines(t *testing.T) {
	result, err := AsciiArt("\n\n", "standard")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != "\n\n" {
		t.Errorf("Expected exactly two newlines, got: %q", result)
	}
}

func TestAsciiArtHandlerInvalidCharacter(t *testing.T) {
	body := strings.NewReader("text=Hello+😊&banner=standard")
	req := httptest.NewRequest("POST", "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	asciiArtHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestTemplatesParse(t *testing.T) {
	templates := []string{
		"templates/home.html",
		"templates/landing.html",
		"templates/legal.html",
		"templates/error.html",
	}
	for _, tmplPath := range templates {
		_, err := template.ParseFiles(tmplPath)
		if err != nil {
			t.Errorf("Failed to parse template file %s: %v", tmplPath, err)
		}
	}
}