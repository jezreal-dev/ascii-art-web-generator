package main

import (
	"os"
	"strings"
	"errors"
	"html/template"
	"net/http"
	"encoding/json"
)

// PageData holds the context variables passed to the HTML templates.
type PageData struct {
	Text     string
	Banner   string
	AsciiArt string
}

// APIRequest defines the JSON payload received from the client.
type APIRequest struct {
	Text   string `json:"text"`
	Banner string `json:"banner"`
}
// APIResponse defines the JSON payload sent back to the client.
type APIResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}


// homeHandler serves the application landing page.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed: Use GET to access the homepage.")
		return
	}

	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "404 Not Found: The page you are looking for does not exist.")
		return
	}

	if _, err := os.Stat("templates/landing.html"); os.IsNotExist(err) {
		renderError(w, http.StatusNotFound, "404 Not Found: Template file (landing.html) is missing.")
		return
	}

	tmpl, err := template.ParseFiles("templates/landing.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "500 Internal Server Error: Failed to parse landing.html template.")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, nil)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "500 Internal Server Error: Failed to render landing.html template.")
		return
	}
}

// generatorHandler renders the form interface for the generator.
func generatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed: Use GET to view the generator.")
		return
	}

	if _, err := os.Stat("templates/home.html"); os.IsNotExist(err) {
		renderError(w, http.StatusNotFound, "404 Not Found: Template file (home.html) is missing.")
		return
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "500 Internal Server Error: Failed to parse home.html template.")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, PageData{Banner: "standard"})
	if err != nil {
		renderError(w, http.StatusInternalServerError, "500 Internal Server Error: Failed to render home.html template.")
		return
	}
}

// asciiArtHandler processes form inputs and renders the resulting ASCII art.
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed, "HTTP Method Not Allowed: Use POST to generate ASCII art.")
		return
	}

	text := r.FormValue("text")
	if text == "" {
		renderError(w, http.StatusBadRequest, "Bad Request: Input text cannot be empty.")
		return
	}

	banner := r.FormValue("banner")
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		renderError(w, http.StatusBadRequest, "Bad Request: Invalid banner style selected.")
		return
	}

	result, err := AsciiArt(text, banner)
	if err != nil {
		if err.Error() == "Invalid character in input" {
			renderError(w, http.StatusBadRequest, "Bad Request: Input contains invalid characters. Only printable ASCII characters (32-126) are allowed.")
		} else if err.Error() == "banner file not found" || err.Error() == "corrupted banner file" {
			renderError(w, http.StatusNotFound, "404 Not Found: The selected banner file is missing or corrupted.")
		} else {
			renderError(w, http.StatusInternalServerError, "500 Internal Server Error: An error occurred while generating ASCII art: "+err.Error())
		}
		return
	}

	if _, err := os.Stat("templates/home.html"); os.IsNotExist(err) {
		renderError(w, http.StatusNotFound, "Template file (home.html) is missing.")
		return
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error: Failed to parse home.html template.")
		return
	}
	w.WriteHeader(http.StatusOK)

	data := PageData{
		Text:     text,
		Banner:   banner,
		AsciiArt: result,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "500 Internal Server Error:\nError rendering template", http.StatusInternalServerError)
		return
	}
}

// renderError renders a customized HTML error page with the given HTTP status.
func renderError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, msg, status)
		return
	}
	tmpl.Execute(w, map[string]interface{}{
		"Status":  status,
		"Message": msg,
	})
}

// AsciiArt maps the input text to its corresponding ASCII art representation using the selected banner font.
func AsciiArt(input string, banners string) (string, error) {
	filePath := "banners/" + banners + ".txt"

	inputFile, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.New("banner file not found")
	}

	content := strings.ReplaceAll(string(inputFile), "\r\n", "\n")
	inputFileLines := strings.Split(content, "\n")

	// Validate standard banner length (95 characters * 9 lines per block = 855 lines)
	if len(inputFileLines) < 855 {
		return "", errors.New("corrupted banner file")
	}

	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\\n", "\n")

	if input == "" {
		return "", nil
	}

	// Direct return if input consists purely of line breaks
	onlyNewLine := true
	for _, char := range input {
		if char != '\n' {
			onlyNewLine = false
			break
		}
	}
	if onlyNewLine {
		return input, nil
	}

	words := strings.Split(input, "\n")
	result := ""

	for _, word := range words {
		if word == "" {
			result += "\n"
			continue
		}

		for i := 0; i < 8; i++ {
			for _, char := range word {
				if char < 32 || char > 126 {
					return "", errors.New("Invalid character in input")
				}
				// Math lookup: Jump to char block, skip empty separator (+1), fetch current row (i)
				result += inputFileLines[i+(int(char-' ')*9)+1]
			}
			result += "\n"
		}
	}
	return result, nil
}

// privacyHandler serves the Privacy Policy page.
func privacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed: Use GET to access the privacy policy page.")
		return
	}

	tmpl, err := template.ParseFiles("templates/legal.html")
	if err != nil {
		renderError(w, http.StatusNotFound, "404 Not Found: Template file (legal.html) is missing.")
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, map[string]string{
		"Title": "Privacy Policy",
		"Type":  "privacy",
	})
}

// termsHandler serves the Terms of Service page.
func termsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed: Use GET to access the terms page.")
		return
	}

	tmpl, err := template.ParseFiles("templates/legal.html")
	if err != nil {
		renderError(w, http.StatusNotFound, "404 Not Found: Template file (legal.html) is missing.")
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, map[string]string{
		"Title": "Terms of Service",
		"Type":  "terms",
	})
}

func apiAsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Error: "405 Method Not Allowed: Use POST."})
		return
	}

	var req APIRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "400 Bad Request: Invalid JSON payload."})
		return
	}

	if req.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "400 Bad Request: Text field cannot be empty."})
		return
	}

	if req.Banner != "standard" && req.Banner != "shadow" && req.Banner != "thinkertoy" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "400 Bad Request: Invalid banner style. Allowed values are 'standard', 'shadow', 'thinkertoy'."})
		return
	}

	result, err := AsciiArt(req.Text, req.Banner)
	if err != nil {
		if err.Error() == "Invalid character in input" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(APIResponse{Error: "400 Bad Request: Input contains invalid characters. Only printable ASCII characters (32-126) are allowed."})
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIResponse{Error: "404 Not Found: The selected banner file is missing or corrupted."})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Result: result})
}