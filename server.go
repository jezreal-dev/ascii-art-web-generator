package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"
	"errors"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
    http.Error(w,"404 Not Found:\nError loading page", http.StatusNotFound)
    return
}

  tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, "404 Not Found:\nError loading template", http.StatusNotFound)
		return
    }

	w.WriteHeader(http.StatusOK)

    err = tmpl.Execute(w, "")
	if err != nil {
    	http.Error(w, "500 Internal Server Error:\nError rendering template", http.StatusInternalServerError)
    	return
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Error:\nMethod not allowed", http.StatusMethodNotAllowed)
        return
	}

	text := r.FormValue("text")
	if text == "" {
        http.Error(w, "400 Bad Request:\nText cannot be empty", http.StatusBadRequest)
        return
    }

	banner := r.FormValue("banner")	
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
    	http.Error(w, "400 Bad Request:\nInvalid banner\nIncorrect Request", http.StatusBadRequest)
    	return
	}
	
	result, err := AsciiArt(text, banner)
	if err != nil {
		http.Error(w, "404 Not Found:\nError loading banner file", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "404 Not Found:\nError loading template", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = tmpl.Execute(w, result)
 	if err != nil {
    	http.Error(w, "500 Internal Server Error:\nError rendering template", http.StatusInternalServerError)
    	return
	}
}


func AsciiArt(input string, banners string) (string, error) {
	filePath := "banners/" + banners + ".txt"

	inputFile, err := os.ReadFile(filePath)
	if err != nil {
	 return "", errors.New("banner file not found")
	}

	content := strings.ReplaceAll(string(inputFile), "\r\n", "\n")
	inputFileLines := strings.Split(content, "\n")

	words := strings.Split(input, "\\n")
	result := ""

	for _, word := range words {
		if input == "" {
    		return "", nil
		}
		if word == "" {
			result += "\n"
			continue
		}
		for i := 0; i < 8; i++ { 
			for _, char := range word {
				result += inputFileLines[i+(int(char-' ')*9)+1]
			}
			result += "\n"
		}
	}
	return result, nil
}