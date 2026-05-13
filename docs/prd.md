# PRD - ASCII ART WEB

## 1. Purpose

Build a Go HTTP server that provides a web GUI for the previous ASCII-art project. Users submit text and a banner choice through a webpage and receive the rendered ASCII-art result.

## 2. Scope

### In scope
- Serve an HTML main page with a text input and banner selector
- Accept form submissions for ASCII-art rendering
- Support `standard`, `shadow`, and `thinkertoy`
- Return correct HTTP status codes
- Use Go HTML templates from the root `templates/` directory
- Provide a more polished, modern, and user-friendly interface than the base web version
- Include CSS-driven styling and interaction feedback

### Out of scope
- Building a JavaScript-heavy frontend
- Using external Go packages
- Editing banner templates

## 3. Users and Use Cases

- A user visits `GET /` and sees the form.
- A user submits text and banner via `POST /ascii-art`.
- The server renders the ASCII art and displays it in HTML.
- An auditor verifies correct responses and status codes for valid and invalid requests.

## 4. HTTP Contract

### 4.1 `GET /`
- Returns the main page.
- Status code: `200 OK`
- The page links to a stylesheet served from `GET /styles.css`.
- The page supports keyboard submission with `Enter`, while `Shift+Enter` keeps multiline input possible.

### 4.2 `POST /ascii-art`
- Accepts text and banner from a form submission.
- On success, returns rendered output in HTML with `200 OK`.
- On invalid input, returns `400 Bad Request` with a human-readable explanation.
- On missing template/banner resource, returns `404 Not Found` with a specific explanation.
- On unexpected server failure, returns `500 Internal Server Error` with a specific explanation.

## 5. Input Rules

- Text input may be empty or non-empty.
- Empty input is accepted and returns a successful HTML response with an empty rendered result.
- Banner input must be one of:
  - `standard`
  - `shadow`
  - `thinkertoy`
- Unsupported banners must not be rendered silently.

## 6. Execution Model

- The project follows a linear web pipeline:
  - HTTP request
  - form parsing and validation
  - banner loading
  - ASCII-art rendering
  - template rendering
  - HTTP response

## 7. Architecture

### HTTP/application layer
- Handles routing, form parsing, validation, status codes, and template responses.

### Banner/render layer
- Loads banner files and converts text into ASCII-art output.
- Supports standard ASCII-art newline behavior for both actual line breaks and `\n` sequences.

### Template layer
- Renders HTML views using server-side Go templates.
- Presents user-facing error feedback inside the page for form-related failures.
- Presents standalone styled error pages for route-level or transport-level failures such as `404 Not Found`.
- Provides browser-side result actions such as copy, save-as-image, and share without changing the server API contract.
- Must keep text readable and visually distinct regardless of background and accent colors.
- Must remain responsive and consistent across screen sizes.

## 8. Testing Strategy

### Unit tests
- Banner loading behavior
- ASCII-art rendering behavior

### HTTP handler tests
- `GET /` success
- `POST /ascii-art` success
- `POST /ascii-art` invalid form/banner handling
- missing route -> `404`
- template/banner failures -> appropriate error code

## 9. Acceptance Criteria

- The server runs in Go.
- `GET /` returns the main page with form controls.
- `POST /ascii-art` renders submitted text with the selected banner.
- Status codes match the subject requirements.
- Templates live under `templates/`.
- The page is visually improved, responsive, and provides clearer feedback than the base web version.
- `gofmt`, `go vet`, and `go test` all pass.
