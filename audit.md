# Friendly Audit Guide

This guide helps an auditor verify the project quickly, one requirement at a time, with concrete commands and expected behavior.

## 1. Start the Server

Run the server:

```bash
go run .
```

Open the website in a browser:

```text
http://localhost:8080
```

Expected:
- The home page loads.
- A text area is visible.
- Banner choices `standard`, `shadow`, and `thinkertoy` are visible.
- The page is styled with CSS.

## 2. Quick Route Checks

Home page:

```bash
curl -i http://localhost:8080/
```

Expected:
- Status `200 OK`
- HTML response

Stylesheet:

```bash
curl -i http://localhost:8080/styles.css
```

Expected:
- Status `200 OK`
- `Content-Type: text/css`

Missing route:

```bash
curl -i http://localhost:8080/missing
```

Expected:
- Status `404 Not Found`
- Message explains that the route does not exist

Wrong method:

```bash
curl -i http://localhost:8080/ascii-art
```

Expected:
- Status `400 Bad Request`
- Message explains that `GET` is not allowed for `/ascii-art`

## 3. Browser Form Checks

Use the web page for the following cases.

### Standard banner: `123??`

Submit:

```text
text   = 123??
banner = standard
```

Expected rendered output:

```text
                     ___    ___
 _   ____    _____  |__ \  |__ \
/ | |___ \  |___ /     ) |    ) |
| |   __) |   |_ \    / /    / /
| |  / __/   ___) |  |_|    |_|
|_| |_____| |____/   (_)    (_)
                                
                                
```

### Standard banner: two-line sample

Submit:

```text
line 1 = {123}
line 2 = <Hello> (World)!
banner = standard
```

Expected:
- The result contains both rendered lines.
- The output remains readable and aligned.

### Shadow banner sample

Submit:

```text
text   = $% \"=
banner = shadow
```

Expected:
- The result matches the subject example.
- The output is readable and does not break the page.

### Thinkertoy banner sample

Submit:

```text
text   = 123 T/fs#R
banner = thinkertoy
```

Expected:
- The result matches the subject example.

## 4. Newline Checks

Actual line break:

```bash
curl -i -X POST http://localhost:8080/ascii-art -d "text=Hello%0AThere&banner=standard"
```

Expected:
- Status `200 OK`
- HTML contains two rendered ASCII-art blocks

Literal `\n` sequence:

```bash
curl -i -X POST http://localhost:8080/ascii-art -d "text=Hello\\nThere&banner=standard"
```

Expected:
- Status `200 OK`
- HTML contains two rendered ASCII-art blocks

## 5. Empty Input Check

```bash
curl -i -X POST http://localhost:8080/ascii-art -d "text=&banner=shadow"
```

Expected:
- Status `200 OK`
- The page stays usable
- The form does not crash
- A neutral ready state is shown instead of a broken result block

## 6. Bad Request Checks

Unsupported banner:

```bash
curl -i -X POST http://localhost:8080/ascii-art -d "text=hello&banner=invalid"
```

Expected:
- Status `400 Bad Request`
- Message explains that the banner is unsupported
- Message suggests the supported values: `standard`, `shadow`, `thinkertoy`

## 7. Template and File Failure Checks

These are mostly covered by automated tests. If the auditor wants to verify manually:

- rename or remove `templates/index.html`
- restart the server or retry the request

Expected:
- Status `404 Not Found`
- Message explains that the required HTML template is missing

If a banner file is missing from `testdata/banners/`, posting with that banner should produce:

- Status `404 Not Found`
- Message explains that the banner resource could not be loaded

## 8. Automated Verification

Run:

```bash
gofmt -w .
go vet ./...
go test ./...
```

Expected:
- All commands complete successfully

## 9. What the Auditor Should Conclude

The project is ready if it demonstrates:
- standard-library-only Go code
- HTML templates in the root `templates/` directory
- working handlers for `GET /` and `POST /ascii-art`
- correct status handling for `200`, `400`, `404`, and `500`
- working ASCII-art rendering for the supported banners
- stable behavior for empty input and newline input
- clear, user-friendly error messages instead of vague responses
