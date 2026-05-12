# ascii-art-web

## Description

`ascii-art-web` is a Go HTTP server that provides a web interface for generating ASCII art from text using selectable banner templates.

## Authors

- ebimai
- gxylinas

## Usage: how to run

```bash
go run .
```

Then open:

```text
http://localhost:8080
```

Submit text through the main page and choose one of the supported banners.
The server also prints the local URL in the terminal so it is easy to open directly.

## Implementation details: algorithm

The server will:

1. Serve the main HTML page on `GET /`
2. Accept submitted text and banner on `POST /ascii-art`
3. Load the selected banner from the filesystem
4. Render the text into ASCII art
5. Return the result inside an HTML template

## Instructions

- The HTTP server is written in Go.
- Only the Go standard library is allowed.
- HTML templates live in the root `templates/` directory.
- The code follows small, focused packages for HTTP, banner loading, and rendering.
- The page styling is served from `GET /styles.css`.
- Supported banners are `standard`, `shadow`, and `thinkertoy`.
- Banner assets are loaded from `testdata/banners/`.
- The HTTP layer depends on a renderer interface so web logic stays separate from banner/render logic.
- Empty text is accepted and returns a successful HTML response.
- Invalid requests return detailed, user-friendly error messages so the browser user understands what went wrong.
- Keyboard shortcut: `Enter` submits the form and `Shift+Enter` inserts a new line inside the text box.
- Rendered output includes browser-side actions for copy, save-as-image, and share, with fallback behavior when a browser does not support a feature.

## Audit Help

- See [audit.md](audit.md) for a friendly step-by-step audit guide with ready-to-run commands.

## Verification

Run:

```bash
gofmt -w .
go vet ./...
go test ./...
```
