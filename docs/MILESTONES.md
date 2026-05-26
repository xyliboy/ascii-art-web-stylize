# Milestones

## Mandatory

### Milestone 1: go.mod
Description: Define the module name and Go version for the project.
Dependencies: None
Acceptance Criteria: go.mod exists in the root directory with the correct module name and Go version.

### Milestone 2: internal/ascii/ascii.go
Description: Implement the ASCII art logic (LoadBanner and Render functions).
Dependencies: Milestone 1
Acceptance Criteria: ascii.go correctly loads banner files and renders ASCII art from text input.

### Milestone 3: internal/handlers/handlers.go
Description: Implement the HTTP handlers for GET / and POST /ascii-art.
Dependencies: Milestone 2
Acceptance Criteria:
- GET / returns 200 and serves the main HTML page
- POST /ascii-art returns 200 and displays the ASCII art result
- POST /ascii-art returns 400 for empty or invalid input
- POST /ascii-art returns 404 for unknown banner
- POST /ascii-art rejects path-like or unsupported banner names
- Returns 500 for unhandled errors

### Milestone 4: templates/index.html
Description: Implement the main HTML page with a form for text input and banner selection.
Dependencies: Milestone 3
Acceptance Criteria:
- Page has a text input field
- Page has a banner selection (default: standard)
- Page has a submit button that sends a POST request to /ascii-art
- Page displays the ASCII art result

### Milestone 5: templates/styles.css
Description: Implement the base styling for the main page.
Dependencies: Milestone 4
Acceptance Criteria: Page is visually clear and readable.

### Milestone 6: README.md
Description: Write the project documentation.
Dependencies: Milestone 1
Acceptance Criteria: README.md contains Description, Authors, Usage and Implementation details sections.

## Stylize

### Milestone 7: Responsive layout
Description: Make the page work correctly on both mobile and desktop screens.
Dependencies: Milestone 5
Acceptance Criteria: Layout adjusts properly when the browser width changes.

### Milestone 8: Consistent design
Description: Apply a consistent color palette, typography, and spacing across the whole page.
Dependencies: Milestone 5
Acceptance Criteria: Every element follows the same visual style.

### Milestone 9: Interactive elements
Description: Add hover and focus states to buttons, inputs, and banner selector. Implement app.js with Copy, Save Image, and Share buttons, typewriter placeholder animation, Enter-to-submit shortcut, and arrow key navigation between banners.
Dependencies: Milestone 5
Acceptance Criteria: The user gets visual feedback when interacting with any element. Copy/Save/Share buttons work on the result. Keyboard navigation works across banners.

### Milestone 10: User feedback
Description: Show the user which banner is selected and what state the form is in.
Dependencies: Milestone 4
Acceptance Criteria: Selected banner is clearly indicated. Errors are displayed with a clear message.

## Tests

### Milestone 11: internal/ascii/ascii_test.go
Description: Implement unit tests for the ASCII art logic (LoadBanner and Render).
Dependencies: Milestone 2
Acceptance Criteria:
- LoadBanner loads all three banners with 95 characters each
- LoadBanner returns an error for invalid or empty banner names
- Render produces 8-row output for a single segment
- Render handles multi-line input, empty input, unsupported characters, and all banners

### Milestone 12: internal/handlers/handlers_test.go
Description: Implement unit and golden tests for the HTTP handlers.
Dependencies: Milestone 3
Acceptance Criteria:
- GET / returns 200
- POST /ascii-art with valid input returns 200
- POST /ascii-art with empty input returns 400
- POST /ascii-art with unknown banner returns 404
- POST /ascii-art with a path-like banner returns 404

### Milestone 13: main_test.go
Description: Implement smoke tests for the HTTP route setup.
Dependencies: Milestone 3
Acceptance Criteria:
- GET / through the router returns 200
- Unknown routes through the router return 404
- POST /ascii-art through the router returns 200 for valid input
- GET /api/banners through the router returns 200

## Bonus

### Milestone 14: GET /api/ascii-art
Description: Implement a REST API endpoint that returns ASCII art as JSON.
Dependencies: Milestone 2
Acceptance Criteria:
- GET /api/ascii-art?text=hello&banner=standard returns 200 with JSON response
- Returns 400 for empty or invalid input
- Returns 404 for unknown banner
- Returns 404 for path-like or unsupported banner names

### Milestone 15: GET /api/banners
Description: Implement a REST API endpoint that returns the list of available banners as JSON.
Dependencies: Milestone 2
Acceptance Criteria: GET /api/banners returns 200 with JSON list of available banners.

### Milestone 16: internal/handlers/api_test.go
Description: Implement unit tests for the REST API endpoints.
Dependencies: Milestone 14, Milestone 15
Acceptance Criteria:
- GET /api/banners returns 200 with JSON list
- GET /api/ascii-art with valid input returns 200 with JSON response
- GET /api/ascii-art with empty text returns 400
- GET /api/ascii-art with unknown banner returns 404
- GET /api/ascii-art with a path-like banner returns 404
