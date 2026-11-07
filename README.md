# Sizzle Calculator

A calculator exposed as a JSON REST API (Go backend).

## Running the backend

```bash
cd backend
go run ./cmd/server
```

Environment variables:

| Variable         | Default                 | Description                 |
|------------------|-------------------------|-----------------------------|
| `PORT`           | `8080`                  | Port the server listens on  |
| `ALLOWED_ORIGIN` | `http://localhost:5173` | Origin allowed by CORS      |

## API manual tests

All POST requests must send `Content-Type: application/json`.

With the server running on `localhost:8080`, these requests were run and produced the responses shown (responses trimmed to status line, relevant headers and body; the `Access-Control-*` headers are identical on every response and `Date` varies).

### Successful operations

```bash
curl -i -X POST localhost:8080/api/v1/add -H 'Content-Type: application/json' -d '{"a":5,"b":3}'
```

```
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Content-Type
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: http://localhost:5173
Content-Type: application/json
Content-Length: 13

{"result":8}
```

```bash
curl -i -X POST localhost:8080/api/v1/evaluate -H 'Content-Type: application/json' -d '{"expression":"2+3*4"}'
```

```
HTTP/1.1 200 OK
Content-Type: application/json

{"result":14}
```

```bash
curl -i -X POST localhost:8080/api/v1/add -H 'Content-Type: application/json' -d '{"a":5,"b":0}'
```

`b: 0` is a valid value (not treated as a missing field):

```
HTTP/1.1 200 OK
Content-Type: application/json

{"result":5}
```

```bash
curl -i localhost:8080/health
```

```
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

### Error cases

Status code rule:

- **400**: the input is malformed (broken JSON, missing field, expression syntax, length).
- **422**: the input is well-formed but the math is undefined (division by zero, negative square root, overflow).
- **413**: the request body exceeds 4 KB.
- **415**: the `Content-Type` header is not `application/json`.

```bash
curl -i -X POST localhost:8080/api/v1/divide -H 'Content-Type: application/json' -d '{"a":5,"b":0}'
```

```
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json

{"error":{"code":"DIVISION_BY_ZERO","message":"cannot divide by zero"}}
```

```bash
curl -i -X POST localhost:8080/api/v1/add -H 'Content-Type: application/json' -d '{"a":5}'
```

```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error":{"code":"MISSING_FIELD","message":"field \"b\" is required"}}
```

```bash
curl -i -X POST localhost:8080/api/v1/add -H 'Content-Type: application/json' -d '{"a":5,"b":3,"c":1}'
```

Unknown fields are rejected:

```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error":{"code":"INVALID_JSON","message":"request body contains an unknown field"}}
```

```bash
curl -i -X POST localhost:8080/api/v1/add -H 'Content-Type: application/json' -d 'basura'
```

```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error":{"code":"INVALID_JSON","message":"request body must be a single valid JSON object"}}
```

```bash
curl -i -X POST localhost:8080/api/v1/evaluate -H 'Content-Type: application/json' -d '{"expression":"2+*3"}'
```

```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error":{"code":"INVALID_EXPRESSION","message":"invalid expression at position 2"}}
```

```bash
curl -i -X POST localhost:8080/api/v1/sqrt -H 'Content-Type: application/json' -d '{"a":-4}'
```

```
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json

{"error":{"code":"NEGATIVE_SQRT","message":"cannot calculate square root of negative number"}}
```

```bash
curl -i -X POST localhost:8080/api/v1/add -d '{"a":5,"b":3}'
```

Missing `Content-Type` header:

```
HTTP/1.1 415 Unsupported Media Type
Content-Type: application/json

{"error":{"code":"UNSUPPORTED_MEDIA_TYPE","message":"Content-Type must be application/json"}}
```

```bash
curl -i localhost:8080/api/v1/add
```

Wrong HTTP method:

```
HTTP/1.1 405 Method Not Allowed
Allow: POST
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff

Method Not Allowed
```

## Known limitations / future work

| Limitation | Why it is not fixed |
|---|---|
| No explicit recursion-depth limit in the parser | The 200-byte input cap already bounds depth to ~200 levels; a depth counter is the next improvement. |
| `len()` counts bytes, not characters | The tokenizer only accepts ASCII, where bytes = characters. The error message says "200 bytes" to be exact. |
| Router 404/405 responses are plain text | Go's mux generates them automatically; returning JSON requires a custom router or wrapping the mux. |
| CORS preflight answers on any path | Harmless for a stateless API without authentication. |
| `writeJSON` does not encode into a buffer first | The only payload is a finite `float64` (guaranteed by `finite()`), so encoding cannot fail. |
| `handleBinary` and `handleSqrt` duplicate decode/validate flow | They are two small handlers; a further abstraction would hurt readability. |
| No `slog` or dependency injection in the router | Out of scope here; it would be the first thing to add in a real service. |
| `ALLOWED_ORIGIN="*"` is not validated | Configuration is the deployer's responsibility; the expected value is documented above (a single origin such as `http://localhost:5173`). |
