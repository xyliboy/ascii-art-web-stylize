# Agent Handoff Prompt

Use the following prompt when another developer wants to brief a different AI agent about the current state of the project.

## Ready-to-use Prompt

```text
You are helping on the repository `ascii-art-web`.

Before suggesting implementation changes, first read and align yourself with these files:
- `projectinfo.txt`
- `AI/AGENTS.md`
- `AI/summary.txt`
- `docs/prd.md`
- `docs/architecture/task-ownership.md`
- `README.md`

Current project state:
- This is a new Go project for the `ascii-art-web` subject.
- The project already has working implementation in the HTTP, banner-loading, and rendering layers.
- The goal is to build an HTTP server with:
  - `GET /`
  - `POST /ascii-art`
- The web UI must allow:
  - text input
  - banner selection (`standard`, `shadow`, `thinkertoy`)
  - form submission for ASCII-art rendering
- Required HTTP status codes:
  - `200 OK`
  - `400 Bad Request`
  - `404 Not Found`
  - `500 Internal Server Error`
- Only the Go standard library is allowed.
- HTML templates must live in the root `templates/` directory.

Architecture decisions already made:
- The project should be structured with:
  - `main.go` at the repository root
  - `internal/http/`
  - `internal/banner/`
  - `internal/render/`
  - `templates/`
  - `testdata/banners/`
- `main.go` should stay minimal.
- HTTP concerns must stay separate from ASCII-art rendering logic.
- Banner/render logic should stay separate from the web layer.
- A repository-level `audit.md` exists as a practical audit guide.

Team split already decided:
- Developer A owns:
  - `internal/banner/`
  - `internal/render/`
  - `testdata/banners/`
  - core unit tests
- Developer B owns:
  - `internal/http/`
  - `templates/`
  - handler tests
  - web route behavior
- Shared later:
  - `main.go`
  - final integration
  - final README cleanup

Important collaboration rules:
- The primary user is still learning, so explain decisions clearly.
- Do not silently redesign the architecture.
- Do not collapse responsibilities into one big file.
- Keep the project aligned with `AI/AGENTS.md` and `docs/prd.md`.
- Assume the repository state is authoritative over any older planning text.

Your first task:
- Summarize the current repository state.
- Confirm whether the planned architecture still matches the actual files.
- Then propose only the next step for the developer's owned area.
```

## Purpose

This handoff prompt is meant to:
- preserve the current architectural decisions
- prevent a new agent from making incorrect assumptions
- keep both developers aligned before implementation begins
