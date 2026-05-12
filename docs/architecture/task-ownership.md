# Task Ownership

## Goal

Split the project into two parallel workstreams with minimal merge conflicts and clear ownership boundaries.

## Recommended Split

### Developer A - Core ASCII-art logic(ebimai)

**Ownership**
- `internal/banner/`
- `internal/render/`
- `testdata/banners/`
- core unit tests for banner loading and rendering

**Responsibilities**
- Load banner files from the filesystem
- Validate banner structure
- Expose glyph lookup behavior
- Render input text into ASCII art
- Handle newline behavior in rendering

**Should not own**
- HTTP handlers
- HTML templates
- response status codes
- form parsing

## Developer B - Web layer(gyxilinas)

**Ownership**
- `internal/http/`
- `templates/`
- HTTP handler tests

**Responsibilities**
- Implement `GET /`
- Implement `POST /ascii-art`
- Parse form data
- Validate request input
- Return correct HTTP status codes
- Render HTML templates with server-side data

**Should not own**
- banner parsing internals
- glyph rendering logic

## Shared Integration

**Shared files after both parts are stable**
- `main.go`
- `README.md`
- final end-to-end verification

**Integration goals**
- wire the HTTP layer to the rendering layer
- ensure templates render the server output correctly
- verify full behavior with tests

## Merge Strategy

1. Developer A finishes banner/render contracts first.
2. Developer B builds handlers against those contracts.
3. Integration happens only after each side passes its own tests.
4. Final cleanup and README update happen together.

## Review Checklist

- Does each developer stay inside their owned area?
- Are responsibilities mixed between HTTP and rendering layers?
- Are tests located close to the layer they validate?
- Are changes grouped logically for easier review?
