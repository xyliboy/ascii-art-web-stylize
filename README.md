# ascii-art-web-stylize

## Description

`ascii-art-web-stylize` is a Go web application that renders ASCII art in the browser using selectable banner templates. It builds on the original web version by improving the visual design, responsiveness, user feedback, and interaction quality.

## Authors

- ebimai

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
8. CSS and client-side interactions improve readability, responsiveness, and user feedback without changing the server contract.

## Stylize Goals Implemented

- More polished and professional visual design
- Better form guidance and feedback
- Clearer output presentation
- Responsive layout for smaller screens
- Result actions for copy, export, and share
