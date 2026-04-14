# ASCII Art Web — Complete Code Walkthrough

This document explains every single line of code in the ascii-art-web project, the logic behind each decision, and how all the files connect together.

---

## Table of Contents

1. [Project Structure Overview](#project-structure-overview)
2. [How the Whole System Works](#how-the-whole-system-works)
3. [main.go — The Entry Point](#maingo--the-entry-point)
4. [server.go — Handlers and Logic](#servergo--handlers-and-logic)
   - [homeHandler](#homehandler)
   - [asciiArtHandler](#asciiarthandler)
   - [AsciiArt Function](#asciiart-function)
5. [templates/home.html — The Frontend](#templateshomehtml--the-frontend)
   - [HTML Structure](#html-structure)
   - [CSS Styling](#css-styling)
   - [Go Template Syntax](#go-template-syntax)
6. [server_test.go — Tests](#server_testgo--tests)
7. [Banner Files](#banner-files)
8. [HTTP Status Codes Used](#http-status-codes-used)

---

## Project Structure Overview

```
ascii-art-web/
├── main.go              # Entry point — registers routes, starts server
├── server.go            # HTTP handlers + ASCII art generation logic
├── server_test.go       # Unit tests for the AsciiArt function
├── go.mod               # Go module definition (created by `go mod init`)
├── README.md            # Project documentation
├── walkthrough.md       # This file
├── templates/
│   └── home.html        # HTML template — the page the user sees
└── banners/
    ├── standard.txt     # Standard ASCII art font
    ├── shadow.txt       # Shadow ASCII art font
    └── thinkertoy.txt   # Thinkertoy ASCII art font
```

**Why this structure?**
- `main.go` and `server.go` are separate to keep concerns separated — `main.go` handles startup, `server.go` handles logic.
- `templates/` holds HTML files — the spec requires this folder name.
- `banners/` holds the font data files — kept separate from templates because they serve a different purpose (data vs. presentation).

---

## How the Whole System Works

Here's the complete flow of a user request:

```
User opens browser → types localhost:8080
        ↓
Browser sends GET request to "/"
        ↓
Go server receives it → homeHandler runs
        ↓
homeHandler loads templates/home.html and sends it back
        ↓
User sees the form → types text, picks a banner, clicks "Generate"
        ↓
Browser sends POST request to "/ascii-art" with form data
        ↓
Go server receives it → asciiArtHandler runs
        ↓
asciiArtHandler extracts text and banner from the form
        ↓
Calls AsciiArt(text, banner) to generate the art
        ↓
Loads templates/home.html again, passes the art result to it
        ↓
Template renders with the ASCII art inside the <pre> tag
        ↓
User sees the ASCII art on the page
```

---

## main.go — The Entry Point

```go
package main
```
**Line 1:** Every Go program starts with a package declaration. `package main` is special — it tells Go this is an executable program, not a library. Go looks for a `main` package to know where to start.

---

```go
import (
    "fmt"
    "log"
    "net/http"
)
```
**Lines 3–7:** Import three standard library packages:
- `fmt` — for printing formatted output to the terminal (used for the startup message)
- `log` — for logging fatal errors (used if the server fails to start)
- `net/http` — Go's built-in HTTP server package (handles all web server functionality)

**Why these three?** `net/http` is the core of any Go web server. `fmt` gives us console output. `log` gives us `log.Fatal` which both prints an error AND stops the program — essential for startup failures.

---

```go
func main() {
```
**Line 9:** The `main` function — this is where Go starts executing. Every Go program must have exactly one `main` function in the `main` package.

---

```go
    http.HandleFunc("/", homeHandler)
```
**Line 10:** Register a route. This tells the Go server: "When any browser request comes to the path `/`, call the function `homeHandler` to handle it."
- `"/"` — the route pattern (the home page URL)
- `homeHandler` — the function to call (defined in `server.go`)

**Important:** In Go, the `"/"` pattern is special — it acts as a **catch-all**. Any URL that doesn't match a more specific route will fall through to this handler. That's why `homeHandler` needs to check if the actual path is exactly `"/"` (we handle this in `server.go`).

---

```go
    http.HandleFunc("/ascii-art", asciiArtHandler)
```
**Line 11:** Register a second route. "When a request comes to `/ascii-art`, call `asciiArtHandler`." This is where the form submission goes.

---

```go
    fmt.Println("Server running at http://localhost:8080")
```
**Line 13:** Print a message to the terminal so the developer knows the server is about to start. This appears when you run `go run .` — it confirms the server is alive.

---

```go
    err := http.ListenAndServe(":8080", nil)
```
**Line 15:** Start the HTTP server.
- `":8080"` — listen on port 8080 on all network interfaces. The colon before the number means "any available address on this machine."
- `nil` — use Go's default request multiplexer (`DefaultServeMux`), which already knows about the routes we registered above with `HandleFunc`.
- `err` — this function returns an error if the server fails to start (e.g., port already in use).

**Why port 8080?** Port 80 is the standard HTTP port but requires admin privileges. Port 8080 is the conventional alternative for development servers.

---

```go
    if err != nil {
        log.Fatal("Error starting server:", err)
    }
```
**Lines 16–18:** Error handling. If `ListenAndServe` returns an error (meaning the server couldn't start):
- `log.Fatal` prints the error message AND immediately terminates the program with exit code 1.
- This is appropriate here because if the server can't start, there's nothing else the program can do.

**Why `log.Fatal` instead of `fmt.Println`?** `fmt.Println` just prints — the program would continue running (doing nothing). `log.Fatal` prints AND exits, which is the correct behavior for an unrecoverable startup error.

---

## server.go — Handlers and Logic

```go
package main
```
**Line 1:** Same package as `main.go`. In Go, all `.go` files in the same directory with `package main` are compiled together as one program. This is why `homeHandler` and `asciiArtHandler` defined here are accessible from `main.go` without importing anything.

---

```go
import (
    "html/template"
    "net/http"
    "os"
    "strings"
)
```
**Lines 3–8:** Four imports:
- `html/template` — Go's HTML template engine. Parses `.html` files and injects data into them. Uses `html/template` (not `text/template`) because it automatically escapes HTML to prevent XSS attacks.
- `net/http` — needed for `http.ResponseWriter`, `http.Request`, `http.Error`, and status code constants.
- `os` — for `os.ReadFile` to read the banner `.txt` files from disk.
- `strings` — for `strings.ReplaceAll` and `strings.Split` to process file content and user input.

---

### homeHandler

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
```
**Line 10:** Handler function for the `GET /` route.
- `w http.ResponseWriter` — the outgoing response. We write data to `w` and it gets sent back to the user's browser.
- `r *http.Request` — the incoming request. Contains the URL, method, form data, headers, etc. The `*` means it's a pointer (reference to the request object).

**Every handler in Go has this exact same signature** — two parameters, `w` and `r`.

---

```go
    if r.URL.Path != "/" {
        http.Error(w, "Error loading page", http.StatusNotFound)
        return
    }
```
**Lines 11–14:** Route guard — prevents the catch-all behavior of `"/"`.

**Why is this needed?** Go's `http.HandleFunc("/", ...)` matches not just `/` but also `/anything`, `/foo/bar`, etc. — any URL that doesn't match another registered route. Without this check, visiting `localhost:8080/nonexistent` would show the home page instead of a 404 error.

- `r.URL.Path` — the actual URL path from the browser request
- `!= "/"` — if the path is anything other than exactly `/`
- `http.Error(w, "Error loading page", http.StatusNotFound)` — sends a plain text error response with HTTP status 404
- `http.StatusNotFound` — Go's constant for the number `404`
- `return` — **critical**: stops the function here. Without `return`, the code below would still execute and try to serve the home page.

---

```go
    tmpl, err := template.ParseFiles("templates/home.html")
```
**Line 16:** Load and parse the HTML template file.
- `template.ParseFiles(...)` reads the file from disk and parses it into a template object.
- Returns two values: `tmpl` (the parsed template) and `err` (any error that occurred).
- `"templates/home.html"` — the file path relative to where the Go program is run from (the project root).

**Why parse every request?** For simplicity. In production, you'd parse once at startup and reuse. For this project, parsing each time is fine and makes the code clearer.

---

```go
    if err != nil {
        http.Error(w, "Error loading template", http.StatusNotFound)
        return
    }
```
**Lines 17–20:** If the template file couldn't be found or parsed:
- Send a 404 error to the browser
- `return` to stop execution

**Why 404?** The spec says "404 Not Found, if nothing is found, for example templates or banners." A missing template file is a "not found" scenario.

**Why not `log.Fatal`?** Because `log.Fatal` would kill the entire server for all users. Inside a handler, we only want to fail for this one request while keeping the server running for everyone else.

---

```go
    tmpl.Execute(w, "")
```
**Line 21:** Render the template and send it to the browser.
- `tmpl.Execute(w, "")` — takes the parsed template, replaces any `{{ . }}` placeholders with the data passed as the second argument, and writes the result to `w` (the response).
- `""` — empty string. On the home page's first load, there's no ASCII art to display yet. We pass an empty string so `{{ . }}` renders as nothing (not `<nil>`).

---

### asciiArtHandler

```go
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
```
**Line 24:** Handler for the `POST /ascii-art` route. Same signature as every Go handler.

---

```go
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
```
**Lines 25–28:** Method guard — only allow POST requests.
- `r.Method` — the HTTP method of the incoming request (GET, POST, PUT, DELETE, etc.)
- `http.MethodPost` — Go's constant for the string `"POST"`
- `!= ` — "not equal to"
- If someone sends a GET request to `/ascii-art` (e.g., by typing it in the browser URL bar), they get a 405 error instead of the handler trying to process a non-existent form.
- `http.StatusMethodNotAllowed` — Go's constant for `405`

**Why is this necessary?** `http.HandleFunc("/ascii-art", asciiArtHandler)` registers the handler for ALL methods, not just POST. Without this guard, a GET request to `/ascii-art` would try to read empty form data and fail confusingly.

---

```go
    text := r.FormValue("text")
    if text == "" {
        http.Error(w, "Text cannot be empty", http.StatusBadRequest)
        return
    }
```
**Lines 30–34:** Extract and validate the user's text input.
- `r.FormValue("text")` — retrieves the value of the form field named `"text"` from the POST request body. The `"text"` matches the `name="text"` attribute on the `<input>` element in `home.html`.
- `text == ""` — if the user submitted the form without typing anything
- `http.StatusBadRequest` — Go's constant for `400`. The spec says "400 Bad Request, for incorrect requests."

---

```go
    banner := r.FormValue("banner")
    if banner == "" {
        http.Error(w, "Banner cannot be empty", http.StatusBadRequest)
        return
    }
```
**Lines 36–40:** Extract and validate the banner selection.
- `r.FormValue("banner")` — gets the selected radio button value. The `"banner"` matches `name="banner"` on the radio inputs in `home.html`.
- If no radio button was selected, `banner` will be an empty string → 400 error.

---

```go
    result := AsciiArt(text, banner)
    if result == "" {
        http.Error(w, "Error Generating Ascii Art", http.StatusInternalServerError)
        return
    }
```
**Lines 42–46:** Generate the ASCII art and validate the output.
- `AsciiArt(text, banner)` — calls the function defined below, passing the user's text and their banner choice.
- `result == ""` — if `AsciiArt` returned an empty string, something went wrong internally (e.g., the banner file couldn't be read).
- `http.StatusInternalServerError` — Go's constant for `500`. The spec says "500 Internal Server Error, for unhandled errors."

---

```go
    tmpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, "Error loading template", http.StatusNotFound)
        return
    }
    tmpl.Execute(w, result)
```
**Lines 48–53:** Load the template and render it with the ASCII art result.
- Same pattern as `homeHandler`, but this time we pass `result` (the generated ASCII art string) instead of `""`.
- When the template renders, `{{ . }}` gets replaced with the ASCII art string.
- The browser receives a complete HTML page with the art already inside the `<pre>` tag.

---

### AsciiArt Function

```go
func AsciiArt(input string, banners string) string {
```
**Line 57:** The core function that generates ASCII art.
- `input string` — the text the user typed (e.g., `"Hello"`)
- `banners string` — the banner name (e.g., `"standard"`, `"shadow"`, or `"thinkertoy"`)
- Returns `string` — the generated ASCII art, or `""` on failure

**This is NOT a handler function** — it doesn't have `w` and `r` parameters. It's a pure helper function that takes input and returns output.

---

```go
    filePath := "banners/" + banners + ".txt"
```
**Line 58:** Build the file path dynamically.
- If `banners` is `"standard"`, this becomes `"banners/standard.txt"`
- If `banners` is `"shadow"`, this becomes `"banners/shadow.txt"`
- String concatenation with `+`

---

```go
    inputFile, err := os.ReadFile(filePath)
    if err != nil {
        return ""
    }
```
**Lines 60–63:** Read the banner file from disk.
- `os.ReadFile(filePath)` — reads the entire file into memory as a `[]byte` (byte slice).
- Returns `inputFile` (the file contents) and `err` (any error).
- If the file doesn't exist or can't be read, return `""` to signal failure.

---

```go
    content := strings.ReplaceAll(string(inputFile), "\r\n", "\n")
```
**Line 65:** Normalize line endings.
- `string(inputFile)` — convert `[]byte` to a `string` (Go requires explicit type conversion).
- `strings.ReplaceAll(..., "\r\n", "\n")` — replace Windows line endings (`\r\n`) with Unix line endings (`\n`).
- **Why?** Banner files might have been created or edited on Windows. Windows uses `\r\n` for newlines, Unix/Mac uses `\n`. If we don't normalize, `strings.Split` would leave invisible `\r` characters in each line, corrupting the ASCII art alignment.

---

```go
    inputFileLines := strings.Split(content, "\n")
```
**Line 66:** Split the banner file into individual lines.
- Turns the entire file content into a **slice (array) of strings**, one per line.
- After this, `inputFileLines[0]` is the first line, `inputFileLines[1]` is the second, etc.
- This lets us access any specific line by its index number — which is how we look up character art.

---

```go
    words := strings.Split(input, "\\n")
```
**Line 68:** Split the user's input on the literal `\n` escape sequence.
- `"\\n"` in Go source code represents the **literal two-character string** `\n` (backslash followed by n).
- **Why `\\n` and not `\n`?** In Go, `"\n"` is a real newline character. But the user types a literal backslash-n in the text box (e.g., `Hello\nWorld`). The browser sends this as the characters `\`, `n` — not an actual newline. So we split on the literal `\n`.
- Result: if the user types `Hello\nWorld`, `words` becomes `["Hello", "World"]`.

---

```go
    result := ""
```
**Line 69:** Initialize an empty string to accumulate the ASCII art output. Each line of art gets appended to this.

---

```go
    for _, word := range words {
```
**Line 71:** Loop through each word (each line of text the user wants rendered).
- `range words` — iterates over the `words` slice.
- `_` — the index (we don't need it, so we discard it with `_`).
- `word` — the current word/line being processed.

---

```go
        if word == "" {
            result += "\n"
            continue
        }
```
**Lines 72–75:** Handle empty lines.
- If the user typed `Hello\n\nWorld`, splitting produces `["Hello", "", "World"]`. The empty string `""` represents a blank line between "Hello" and "World".
- `result += "\n"` — add a blank line to the output.
- `continue` — skip to the next word (don't try to render an empty word through the character loop).

---

```go
        for i := 0; i < 8; i++ {
```
**Line 76:** Inner loop — iterate 8 times, once for each line of the ASCII character art.
- Every character in the banner file is represented by exactly **8 lines** of text art.
- `i` is the current line number (0 through 7) of the character being built.

---

```go
            for _, char := range word {
```
**Line 77:** Innermost loop — iterate through each character in the current word.
- `range word` iterates over the string character by character.
- `char` is a `rune` (Go's type for a single character, representing its Unicode/ASCII value).

---

```go
                result += inputFileLines[i+(int(char-' ')*9)+1]
```
**Line 78:** **The core formula** — this is the most important line in the entire project. It looks up the correct line in the banner file for the current character and current row.

Breaking it down piece by piece:

1. **`char - ' '`** — Calculate the character's position relative to the space character.
   - `' '` (space) has ASCII value 32.
   - `'A'` has ASCII value 65, so `'A' - ' '` = 65 - 32 = **33**.
   - This gives us which character block to look at in the banner file.

2. **`(char - ' ') * 9`** — Each character occupies **9 lines** in the banner file (8 lines of art + 1 blank separator line between characters). Multiplying by 9 jumps to the start of the correct character's block.

3. **`+ 1`** — Skip the blank separator line at the start of each character block. The first line of each block is empty, so `+1` moves past it.

4. **`+ i`** — Add the current row number (0–7) to get the specific line within that character's 8-line block.

5. **`int(...)`** — Convert the result to an `int` because Go requires array indices to be integers.

**Example:** To get line 3 (i=2) of the character `'A'` (position 33):
```
Index = 2 + (33 * 9) + 1 = 2 + 297 + 1 = 300
```
So `inputFileLines[300]` gives us the 3rd line of the letter A's ASCII art.

---

```go
            result += "\n"
```
**Line 80:** After printing one row of all characters in the word, add a newline to move to the next row.

---

```go
    return result
```
**Line 83:** Return the complete ASCII art string.

---

## templates/home.html — The Frontend

### HTML Structure

```html
<!DOCTYPE html>
```
**Line 1:** Document type declaration. Tells the browser this is an HTML5 document. Every modern HTML page starts with this.

---

```html
<html lang="en">
```
**Line 2:** Root HTML element. `lang="en"` tells the browser and search engines the page is in English.

---

```html
<head>
    <meta charset="UTF-8">
    <title>ASCII Art Generator</title>
```
**Lines 3–5:**
- `<head>` — contains metadata about the page (not visible content).
- `<meta charset="UTF-8">` — character encoding. UTF-8 supports all characters including special symbols the user might type.
- `<title>` — the text shown in the browser tab.

---

### CSS Styling

```css
body {
    font-family: Arial, sans-serif;
    background-color: #f4f4f4;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px;
}
```
**Lines 7–14:** Style the page body.
- `font-family: Arial, sans-serif` — use Arial font, fall back to any sans-serif font.
- `background-color: #f4f4f4` — light grey background.
- `display: flex` — use CSS Flexbox layout, which makes centering easy.
- `flex-direction: column` — stack child elements vertically (top to bottom).
- `align-items: center` — center everything horizontally on the page.
- `padding: 40px` — space between the page edge and the content.

---

```css
h1 {
    margin-bottom: 20px;
    color: #333;
}
```
**Lines 16–19:** Style the title heading with dark grey color and space below it.

---

```css
form {
    background: #fff;
    padding: 24px;
    border-radius: 8px;
    box-shadow: 0 2px 6px rgba(0,0,0,0.1);
    display: flex;
    flex-direction: column;
    gap: 16px;
    width: 400px;
}
```
**Lines 21–30:** The form appears as a white card.
- `background: #fff` — white background.
- `border-radius: 8px` — rounded corners.
- `box-shadow: 0 2px 6px rgba(0,0,0,0.1)` — subtle shadow beneath the card for depth. `rgba(0,0,0,0.1)` is black at 10% opacity.
- `display: flex; flex-direction: column` — stack form elements vertically.
- `gap: 16px` — equal spacing between each form element.
- `width: 400px` — fixed width.

---

```css
label {
    font-weight: bold;
    color: #444;
}
```
**Lines 32–35:** Make labels bold and dark grey.

---

```css
input[type="text"] {
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 16px;
    width: 100%;
}
```
**Lines 37–43:** Style the text input field.
- `input[type="text"]` — CSS selector that targets only text inputs (not radio buttons or submit buttons).
- `width: 100%` — fill the full width of the form card.

---

```css
.banner-options {
    display: flex;
    gap: 16px;
}
```
**Lines 45–48:** The radio buttons are wrapped in a `<div class="banner-options">`. This makes them display side by side (horizontally) with 16px spacing.

---

```css
input[type="submit"] {
    padding: 10px;
    background-color: #333;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 16px;
    cursor: pointer;
}

input[type="submit"]:hover {
    background-color: #555;
}
```
**Lines 50–62:** Submit button styling.
- Dark background, white text, no default border.
- `cursor: pointer` — changes the mouse cursor to a hand when hovering.
- `:hover` — when the user hovers over the button, it lightens to `#555` for visual feedback.

---

```css
pre {
    margin-top: 30px;
    background: #fff;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 6px rgba(0,0,0,0.1);
    font-size: 14px;
    overflow-x: auto;
    max-width: 90%;
}
```
**Lines 64–73:** Style for the ASCII art output area.
- `overflow-x: auto` — if the ASCII art is wider than the screen, add a horizontal scrollbar instead of breaking the layout.
- `max-width: 90%` — prevent the output from touching the page edges.
- Same white card style as the form for visual consistency.

---

### HTML Body Content

```html
<h1>ASCII Art Generator</h1>
```
**Line 77:** Page title visible to the user.

---

```html
<form action="/ascii-art" method="POST">
```
**Line 79:** The form element that wraps all input fields.
- `action="/ascii-art"` — when submitted, send the data to the `/ascii-art` route (which maps to `asciiArtHandler` in Go).
- `method="POST"` — use the POST HTTP method (data goes in the request body, not the URL).

**Why POST and not GET?**
- GET puts data in the URL (e.g., `?text=Hello&banner=standard`), which has length limits and exposes data.
- POST sends data in the request body — better for form submissions with user input.

---

```html
<label>Enter text:</label>
<input type="text" name="text" placeholder="Type something...">
```
**Lines 80–81:**
- `<label>` — descriptive text telling the user what to type.
- `<input type="text">` — a single-line text input box.
- `name="text"` — **this is critical**. This `name` must match exactly what the Go handler reads with `r.FormValue("text")`. If these don't match, the server won't receive the data.
- `placeholder="Type something..."` — grey hint text that disappears when the user starts typing.

---

```html
<label>Select banner:</label>
<div class="banner-options">
    <label><input type="radio" name="banner" value="standard"> Standard</label>
    <label><input type="radio" name="banner" value="shadow"> Shadow</label>
    <label><input type="radio" name="banner" value="thinkertoy"> Thinkertoy</label>
</div>
```
**Lines 83–88:** Banner selection using radio buttons.
- `type="radio"` — circular buttons where only one can be selected at a time.
- `name="banner"` — all three share the same `name`. This is how HTML knows they're a group (selecting one deselects the others). Matches `r.FormValue("banner")` in Go.
- `value="standard"` / `"shadow"` / `"thinkertoy"` — the value sent to the server when that option is selected. These match exactly what `AsciiArt` expects as the `banners` parameter to build the file path: `"banners/" + "standard" + ".txt"`.
- Each `<input>` is wrapped in a `<label>` so clicking the text also selects the radio button.

---

```html
<input type="submit" value="Generate">
```
**Line 90:** The submit button.
- `type="submit"` — clicking this triggers the form submission.
- `value="Generate"` — the text displayed on the button.

---

```html
</form>

<pre>{{ . }}</pre>
```
**Lines 91–93:** Close the form, then display the result.
- `<pre>` — preformatted text element. It preserves all spaces and line breaks exactly as they appear. Without `<pre>`, the browser would collapse multiple spaces into one and the ASCII art would look broken.
- `{{ . }}` — **Go template syntax**. The dot `.` represents the data passed in `tmpl.Execute(w, data)`. 
  - When `homeHandler` calls `tmpl.Execute(w, "")`, `{{ . }}` becomes an empty string — nothing visible.
  - When `asciiArtHandler` calls `tmpl.Execute(w, result)`, `{{ . }}` becomes the full ASCII art output.

---

## server_test.go — Tests

```go
package main
```
**Line 1:** Same package as the code being tested. In Go, test files in the same package can access all functions directly.

---

```go
import (
    "strings"
    "testing"
)
```
**Lines 3–6:**
- `strings` — used for `strings.Split` and `strings.TrimRight` to analyze output.
- `testing` — Go's built-in testing package. Provides the `*testing.T` type for test assertions.

---

**Test functions** — each follows Go's convention: function name starts with `Test`, takes `*testing.T` as the only parameter.

| Test | What it checks |
|------|-------|
| `TestAsciiArtValidBanner` | `"Hello"` + `"standard"` produces non-empty output |
| `TestAsciiArtShadowBanner` | `"Hello"` + `"shadow"` produces non-empty output |
| `TestAsciiArtThinkertoyBanner` | `"Hello"` + `"thinkertoy"` produces non-empty output |
| `TestAsciiArtInvalidBanner` | `"Hello"` + `"nonexistent"` returns `""` (graceful failure) |
| `TestAsciiArtEmptyInput` | `""` + `"standard"` returns `""` (nothing to render) |
| `TestAsciiArtNewline` | `"Hi\nHi"` produces more than 8 lines (multi-line works) |
| `TestAsciiArtOutputHasEightLinesPerWord` | `"Hi"` produces exactly 8 lines (correct character height) |

Run all tests with:
```bash
go test ./... -v
```

---

## Banner Files

Each banner file (`standard.txt`, `shadow.txt`, `thinkertoy.txt`) follows the same structure:

- Contains the ASCII art for every printable character from space (ASCII 32) to tilde (ASCII 126).
- Each character is represented by **8 lines** of art.
- Between each character block there is **1 blank separator line**.
- So each character occupies **9 lines** total in the file.

**Example structure** (simplified):
```
                        ← blank separator (line 0)
 _                      ← line 1 of '!' (8 lines follow)
| |                     ← line 2 of '!'
| |                     ← line 3 of '!'
| |                     ← ...
|_|
(_)

                        ← blank separator
                        ← line 1 of '"' ...
```

**The formula `i + (int(char-' ') * 9) + 1` navigates this structure:**
- `(char - ' ')` → which character (0 for space, 1 for !, 2 for ", ...)
- `* 9` → jump to that character's block
- `+ 1` → skip the blank separator at the start
- `+ i` → select the specific line (0-7) within the block

---

## HTTP Status Codes Used

| Code | Constant | Where Used | When |
|------|----------|------------|------|
| 200 | (default) | `tmpl.Execute` | Template renders successfully — Go sends 200 automatically |
| 400 | `http.StatusBadRequest` | `asciiArtHandler` | Text or banner field is empty |
| 404 | `http.StatusNotFound` | `homeHandler`, `asciiArtHandler` | Invalid URL path or template file missing |
| 405 | `http.StatusMethodNotAllowed` | `asciiArtHandler` | Request method is not POST |
| 500 | `http.StatusInternalServerError` | `asciiArtHandler` | `AsciiArt` function returns empty (banner file unreadable) |
