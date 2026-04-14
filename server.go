package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
    http.Error(w,"Error loading page", http.StatusNotFound)
    return
}

  tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, "Error loading template", http.StatusNotFound)
		return
    }
    tmpl.Execute(w, "")
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
	}

	text := r.FormValue("text")
	if text == "" {
        http.Error(w, "Text cannot be empty", http.StatusBadRequest)
        return
    }

	banner := r.FormValue("banner")	
	if banner == "" {
        http.Error(w, "Banner cannot be empty", http.StatusBadRequest)
        return
    }
	
	result := AsciiArt(text, banner)
	if result == "" {
		http.Error(w, "Error Generating Ascii Art", http.StatusInternalServerError)
		return
	}	

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusNotFound)
		return
	}
	tmpl.Execute(w, result)
}


func AsciiArt(input string, banners string) string {
	filePath := "banners/" + banners + ".txt"

	inputFile, err := os.ReadFile(filePath)
	if err != nil {
	 return ""
	}

	content := strings.ReplaceAll(string(inputFile), "\r\n", "\n")
	inputFileLines := strings.Split(content, "\n")

	words := strings.Split(input, "\\n")
	result := "" //

	for _, word := range words {
		if word == "" {
			result += "\n"
			continue
		}
		for i := 0; i < 8; i++ { // i
			for _, char := range word {
				result += inputFileLines[i+(int(char-' ')*9)+1]
			}
			result += "\n"
		}
	}
	return result
}