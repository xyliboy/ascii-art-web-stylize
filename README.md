# ascii-art-web-stylize

## Description

ascii-art-web-stylize is a Go web application that renders ASCII art in the browser. The user types a word or short message, picks a banner style, and the server returns the text as ASCII art. The project builds on ascii-art-web by adding a polished, responsive, and user-friendly design with CSS.

## Authors

- gtzimoka
- ebimai
- gxylinas

## Usage: how to run

```bash
go run .
```

Then open in your browser:

```
http://localhost:8080
```

## Implementation details: algorithm

1. The server starts on port 8080.
2. GET / renders the main HTML page from templates/index.html.
3. The user types text, selects a banner, and submits the form.
4. POST /ascii-art reads the text and banner from the form.
5. The banner value is checked against the allowed banners: standard, shadow, and thinkertoy.
6. LoadBanner reads the banner file from banners/ and builds a character map.
7. Render converts each character to its 8-row ASCII art block and joins them.
8. The result is injected into the HTML template and returned to the browser.
9. Errors return the appropriate HTTP status code (400, 404, 500) with a message shown on the page.
10. app.js adds Copy, Save Image (PNG), and Share buttons on the result area, a typewriter animation on the textarea placeholder, Enter-to-submit, and arrow key navigation between banners.
11. Bonus REST API: GET /api/ascii-art?text=hello&banner=standard returns JSON. GET /api/banners returns the list of available banners.

## Tests

Run all tests with:

```bash
go test ./...
```

The test suite covers route setup, HTTP handlers, API endpoints, banner loading, ASCII rendering, and golden output cases.
