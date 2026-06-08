package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"
	"errors"
)


// PageData is a struct that holds the data to be passed to the HTML templates. It includes the original text input by the user, the selected banner, and the generated ASCII art. This struct allows us to easily pass all necessary information to the templates for rendering the home page and displaying the results.
type PageData struct {
	Text     string
	Banner   string
	AsciiArt string
}

// The homeHandler serves the home page and renders the form for user input. It also handles the form submission and generates the ASCII art based on the user's input and selected banner.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// for any method other than GET, return 405 Method Not Allowed
	if r.Method != http.MethodGet {
		http.Error(w, "405 Error: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// if the path is not "/", return 404 Not Found
	if r.URL.Path != "/" {
    	renderError(w, http.StatusNotFound, "404 Not Found: The page you are looking for does not exist.")
    	return
	}

	// Load the template and render the home page with the default banner selected
  	tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        renderError(w, http.StatusNotFound, "404 Not Found: Template file (home.html) is missing.")
		return
    }

	w.WriteHeader(http.StatusOK)
	// we are going to pass PageData with the default banner selected
    err = tmpl.Execute(w, PageData{Banner: "standard"})
	if err != nil {
    	renderError(w, http.StatusInternalServerError, "500 Internal Server Error: Failed to render home.html template.")
    	return
	}
}


// The asciiArtHandler handles the form submission from the home page. It validates the input, generates the ASCII art using the AsciiArt function, and renders the result in the template. It also handles various error cases such as invalid method, empty text, invalid banner selection, and errors during ASCII art generation or template rendering.
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// for any method other than POST, return 405 Method Not Allowed 
	if r.Method != http.MethodPost { 
		renderError(w, http.StatusMethodNotAllowed, "HTTP Method Not Allowed: Use POST to generate ASCII art.")
        return
	}

	// if text is empty, return 400 Bad Request with a message indicating that text cannot be empty
	text := r.FormValue("text")
	if text == "" { 
        renderError(w, http.StatusBadRequest, "Bad Request: Input text cannot be empty.")
        return
    }

	// Validate banner value and return 400 Bad Request if it's invalid 
	banner := r.FormValue("banner")	
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
    	renderError(w, http.StatusBadRequest, "Bad Request: Invalid banner style selected.")
    	return
	}
	
	// Generate ASCII art using the AsciiArt function 
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

	// Load the template and render the result in the template, passing the original text and banner as well
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Internal Server Error: Failed to parse home.html template.")
		return
	}
	w.WriteHeader(http.StatusOK)

	// we are going to pass PageData with the original text, banner, and the generated ASCII art
	data := PageData{
		Text:     text,
		Banner:   banner,
		AsciiArt: result,
	}

	// Render the template with the data
	err = tmpl.Execute(w, data)
 	if err != nil {
    	http.Error(w, "500 Internal Server Error:\nError rendering template", http.StatusInternalServerError)
    	return
	}
}


// renderError is a helper function to render error messages using the error.html template. It takes the response writer, status code, and error message as parameters. If the template cannot be loaded, it falls back to sending a plain text error response.
func renderError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// Fallback to plain text if the error template itself is missing or broken
		http.Error(w, msg, status)
		return
	}
	// Render the error template with the provided status and message
	tmpl.Execute(w, map[string]interface{}{
		"Status": status,
		"Message": msg,
	})
}


// AsciiArt is a function that takes an input string and a banner name, reads the corresponding banner file, and generates the ASCII art representation of the input text based on the banner. It handles various edge cases such as empty input, input with only new lines, and invalid characters. It also ensures that the banner file is properly formatted and contains enough lines to represent all printable ASCII characters.
func AsciiArt(input string, banners string) (string, error) {
	filePath := "banners/" + banners + ".txt"

	inputFile, err := os.ReadFile(filePath)
	if err != nil {
	 return "", errors.New("banner file not found")
	}

	content := strings.ReplaceAll(string(inputFile), "\r\n", "\n")
	inputFileLines := strings.Split(content, "\n")

	// Verify the banner file has enough lines (95 characters * 9 lines = 855 lines)
	if len(inputFileLines) < 855 {
		return "", errors.New("corrupted banner file")
	}

	// Normalize all newlines: replace CRLF and escaped \n with actual LF \n
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\\n", "\n")

	// For empty input
	if input == "" {
		return "", nil
	}

	// If input consists of only new lines 
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
				// Validate character range (printable ASCII characters: 32 to 126)
				if char < 32 || char > 126 {
					return "", errors.New("Invalid character in input")
				}
				result += inputFileLines[i+(int(char-' ')*9)+1]
			}
			result += "\n"
		}
	}
	return result, nil
}