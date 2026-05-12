This document defines the repository-specific development rules for `ascii-art-web`.

## 1. Core Rules

### 1.0 Collaboration Mode
- This project is being developed by two developers working in the same repository.
- Work should be split into clear areas of ownership to reduce conflicts and improve review quality.
- Changes should be grouped logically so each contributor can explain their part of the work.
- The AI assistant must support collaboration, not replace it.

### 1.0.1 Apprentice Context
- The primary user is still learning and should be guided as an apprentice developer.
- Explanations should be clear, practical, and architecture-aware.
- The AI assistant should not silently do all the thinking and implementation without explanation.
- Prefer helping the user understand why a design or code decision is correct.

### 1.1 Standard Library Only
- Use only the Go standard library.
- No external packages are allowed unless the subject explicitly allows them.

### 1.2 Code Style
- Run `gofmt -w .` before considering work complete.
- Functions should stay small, focused, and readable.
- Prefer single-purpose helpers over long multi-responsibility functions.

### 1.3 Naming
- Variables: descriptive camelCase
- Functions: verb-based names
- Types and constants: PascalCase

## 2. Functional Contract

### 2.1 HTTP Endpoints
- `GET /`
  - Returns the main HTML page.
- `POST /ascii-art`
  - Accepts form data with text and banner.
  - Renders ASCII art and returns the result in HTML.

### 2.2 Supported Banners
- `standard`
- `shadow`
- `thinkertoy`

### 2.3 HTTP Status Codes
- `200 OK` for successful requests
- `400 Bad Request` for invalid form submissions or unsupported request formats
- `404 Not Found` for missing routes, templates, or banner files
- `500 Internal Server Error` for unexpected server-side failures

## 3. Architecture Expectations

### 3.1 HTTP Layer
- Owns request routing, request validation, response status codes, and template rendering.
- Must remain separate from ASCII-art rendering logic.

### 3.2 Rendering Layer
- Owns banner loading and ASCII-art generation.
- Should stay independent from HTTP details.

### 3.3 Templates
- HTML templates must live in the root `templates/` directory, per subject requirement.

## 4. Testing Expectations

### 4.1 TDD
- Prefer writing or updating tests before changing behavior.

### 4.2 Unit Tests
- Rendering and banner logic should be tested in isolation.
- HTTP handlers should be tested with `net/http/httptest`.

### 4.3 Behavior Coverage
- Validate successful home page rendering
- Validate successful POST rendering
- Validate invalid banner handling
- Validate empty input handling
- Validate missing routes/templates/banners with proper HTTP status codes

## 5. Tooling and Verification

Before considering work complete, run:

```powershell
gofmt -w .
go vet ./...
go test ./...
```

## 6. AI Collaboration Boundaries

- The AI assistant may help design the architecture, explain concepts, draft code, and review changes.
- The AI assistant must not replace the developers' understanding of the project.
- The AI assistant should avoid taking over the whole project without teaching what is being done.
- Prefer pair-programming style guidance when decisions affect learning or architecture.
- When work can be split safely, prefer assigning clear ownership between the two developers.

## 7. Documentation Alignment

When requirements or implementation decisions change, keep these files aligned:
- `AI/AGENTS.md`
- `AI/summary.txt`
- `docs/prd.md`
- `README.md`
- `projectinfo.txt`
