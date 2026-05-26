# PRD — ascii-art-web-stylize

## Objective

Transform the functional ascii-art-web application into a polished, responsive, and user-friendly web app using CSS. The server logic and ASCII rendering remain unchanged. This project focuses entirely on the UI/UX layer.

## Pages

Single page application:
- `GET /` — main page with form and result area
- `POST /ascii-art` — handles form submission, returns result on same page

## UI Components

1. **Header** — application title
2. **Textarea** — input field for the user's text, with typewriter placeholder animation
3. **Banner selector** — radio buttons for standard / shadow / thinkertoy, with keyboard navigation (arrow keys)
4. **Submit button** — triggers POST to `/ascii-art`, also triggered by Enter key
5. **Result area** — displays ASCII art output inside a `<pre>` element
6. **Error area** — displays error messages clearly when something goes wrong
7. **Action buttons** — Copy (clipboard), Save Image (PNG via Canvas API), Share (Web Share API)

## CSS Requirements

- Responsive layout that works on both mobile and desktop
- Consistent colors and typography throughout the page
- User feedback — clearly show which banner is selected
- Text must remain readable regardless of background or foreground colors chosen
- Interactive elements (buttons, inputs) must have visible hover and focus states

## HTTP Status Codes

- 200 OK — everything worked
- 400 Bad Request — empty input or input with no supported ASCII characters
- 404 Not Found — missing template or banner file
- 500 Internal Server Error — unhandled server error

## Architecture

```
ascii-art-web-stylize/
├── main.go                  — entry point, starts HTTP server on :8080
├── main_test.go             — route setup smoke tests
├── go.mod                   — module definition
├── README.md                — project documentation
├── internal/
│   ├── ascii/
│   │   ├── ascii.go         — LoadBanner() and Render() logic
│   │   └── ascii_test.go    — unit tests for ascii logic
│   └── handlers/
│       ├── handlers.go      — HomeHandler and AsciiArtHandler
│       ├── banners.go       — allowed banner validation
│       ├── handlers_test.go — handler tests
│       ├── api.go           — REST API endpoints (bonus)
│       └── api_test.go      — REST API tests (bonus)
├── templates/
│   ├── index.html           — main HTML page
│   ├── error.html           — custom error page
│   ├── styles.css           — stylesheet
│   └── app.js               — Copy / Save Image / Share buttons, typewriter animation, keyboard navigation
├── banners/
│   ├── standard.txt         — standard banner file
│   ├── shadow.txt           — shadow banner file
│   └── thinkertoy.txt       — thinkertoy banner file
└── docs/
    ├── prd.md               — product requirements
    ├── MILESTONES.md        — project milestones
    └── golden_tests.md      — expected output for audit test cases
```

### Request Flow

```
Browser
  └── GET /
        └── handlers.HomeHandler → renders templates/index.html → 200 OK

Browser
  └── POST /ascii-art (text + banner)
        └── handlers.AsciiArtHandler
              └── validates banner name
              └── ascii.LoadBanner(banner)
                    └── reads banners/<banner>.txt
              └── ascii.Render(text, bannerMap)
              └── renders templates/index.html with result → 200 OK
```

## Constraints

- Only standard Go packages allowed
- HTML templates must live in the `templates/` directory
- Banner files must live in `banners/`
- Code must follow good practices
