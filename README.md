# ascii-art-web-stylize

## Description

`ascii-art-web-stylize` is a Go web application that renders ASCII art in the browser using selectable banner templates. It builds on the original web version by improving the visual design, responsiveness, readability, user feedback, and overall interaction quality.

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

## Implementation details: algorithm

1. The server starts an HTTP application on port `8080`.
2. `GET /` renders the main HTML page.
3. The user submits text and a banner through the form.
4. `POST /ascii-art` validates the request input.
5. The banner loader reads the requested banner file from `testdata/banners/`.
6. The render service converts the submitted text into ASCII art.
7. The server injects the result back into the HTML template.
8. CSS and light browser-side interactions improve readability, responsiveness, scrolling, and feedback without changing the server contract.

## Stylize Goals Implemented

- More polished and professional visual design
- Native-feeling blue page scrollbar instead of a custom page-scroll widget
- More readable galaxy-style animated background
- Consistent visual behavior across the main page and standalone error pages
- Better form guidance and feedback
- Clearer output presentation
- Responsive layout for smaller screens
- Result actions for copy, export, and share

## Audit Help

- Review [audit-questions.txt](/c:/Users/ligor/ascii-art-web-stylize/audit-questions.txt) before the final audit pass.

## Verification

Run:

```bash
gofmt -w .
go vet ./...
go test ./...
```
