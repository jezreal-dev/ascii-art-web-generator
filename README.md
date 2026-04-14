# ASCII Art Web

## Description

ASCII Art Web is a web-based GUI version of the ascii-art project. It runs an HTTP server that allows users to convert text into ASCII art through a browser interface. Users can type any text, choose from three banner styles (standard, shadow, thinkertoy), and view the rendered ASCII art on the same page.

## Authors

- jmomoh

## Usage

### How to run

1. Clone or navigate to the project directory.
2. Make sure you have Go installed (version 1.18+).
3. Run the server:

```bash
go run .
```

4. Open your browser and visit:

```
http://localhost:8080
```

5. Type your text, select a banner, and click **Generate**.

### Example

Input text: `Hello`
Banner: `standard`

Expected output:

```
 _   _          _   _
| | | |        | | | |
| |_| |   ___  | | | |   ___
|  _  |  / _ \ | | | |  / _ \
| | | | |  __/ | | | | | (_) |
|_| |_|  \___| |_| |_|  \___/
```

### Project Structure

```
ascii-art-web/
├── main.go              # Entry point — registers routes and starts the server
├── server.go            # HTTP handlers and ASCII art generation logic
├── server_test.go       # Unit tests for the AsciiArt function
├── go.mod               # Go module definition
├── README.md            # This file
├── templates/
│   └── home.html        # HTML template for the main page
└── banners/
    ├── standard.txt     # Standard banner font
    ├── shadow.txt       # Shadow banner font
    └── thinkertoy.txt   # Thinkertoy banner font
```

## Implementation Details: Algorithm

The ASCII art generation works as follows:

1. **Banner file loading** — Each banner (`standard`, `shadow`, `thinkertoy`) is stored as a `.txt` file in the `banners/` directory. The file contains the ASCII art representation of every printable character (from space onward), each rendered across 8 lines, with a blank separator line between characters.

2. **Input processing** — The user's input is split on the literal `\n` sequence to support multi-line output.

3. **Character lookup** — For each character in the input, its position in the banner file is calculated using:
   ```
   position = (ASCII value of character - ASCII value of space) × 9 + 1
   ```
   Each character occupies 9 lines (8 lines of art + 1 blank separator). The `+1` skips the leading blank line.

4. **Line-by-line rendering** — The output is built by iterating 8 times per word (once per art line), collecting the corresponding banner line for each character, and joining them with newlines.

5. **Template rendering** — The result is passed to Go's `html/template` engine and rendered inside a `<pre>` tag in `templates/home.html` to preserve spacing and formatting.

### HTTP Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/` | Serves the main page with the input form |
| POST | `/ascii-art` | Receives form data, generates ASCII art, returns result |

### HTTP Status Codes

| Code | When |
|------|------|
| 200 OK | Successful request |
| 400 Bad Request | Missing text or banner selection |
| 404 Not Found | Invalid route or missing template |
| 405 Method Not Allowed | Wrong HTTP method used |
| 500 Internal Server Error | ASCII art generation failed |
