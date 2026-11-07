# Repository structure

calculator/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/calculator/   
│   ├── internal/api/          
│   ├── go.mod
│   └── Dockerfile
├── frontend/                  
│   └── Dockerfile
├── docker-compose.yml
├── README.md
└── PROMPTS.md

# API Design
## Separate Endpoints to can tested independently and without server

POST /api/v1/add
POST /api/v1/subtract
POST /api/v1/multiply
POST /api/v1/divide
POST /api/v1/power
POST /api/v1/sqrt
POST /api/v1/percentage
GET /health

## Calculator package (`internal/calculator`)

- Pure logic, no HTTP/JSON dependencies → testable in isolation.
- All operations share the signature `func(...float64) (float64, error)`.
- Domain errors are package-level vars, so callers use `errors.Is`.
- `finite()` rejects NaN/±Inf results (overflow, `Pow(-8, 0.5)`, `Inf/Inf`).
- `Percentage(value, pct)` = pct% of value. Divides before multiplying
  (`value * (pct/100)`) to avoid intermediate overflow.

## Expression evaluator

- Hand-written tokenizer + recursive descent parser. No `eval`, no libraries.
- Grammar (lower rule = higher precedence):
    expression = term   { ("+" | "-") term }
    term       = unary  { ("*" | "/") unary }
    unary      = "-" unary | power
    power      = primary [ "^" unary ]
    primary    = number | "(" expression ")" | "sqrt" "(" expression ")"
- `-2^2 = -4` and `2^3^2 = 512` (right-associative power), matching math convention.
- Reuses the calculator functions, so domain errors (division by zero, etc.)
  behave identically in both endpoints.
- Errors include the byte position: "invalid expression at position 4".
- Scope limits (intentional): ASCII only, no implicit multiplication (`2(3)`),
  `%` not supported inside expressions (ambiguous: modulo vs percent).
- Max expression length enforced in the HTTP handler, which also bounds
  recursion depth for inputs like `((((((...`.

# Answers format
//success
{"result": 42}
//error
{"error": {"code": "DIVISION_BY_ZERO", "message": "Cannot divide by zero"}}

On percentage a is the value and b is the percentage (Remember documentate on readme)

## Testing

- `cmd/server` intentionally has no tests (0% coverage): `main.go` only reads
  env vars and wires `Logger`/`Recover`/`CORS`/`NewRouter` into an `http.Server`.
  All the logic lives in `internal/` and is tested there; a test for `main`
  would only exercise `ListenAndServe` and exist to raise the number.
- Router tests are table-driven (`handlers_test.go`) and run through
  `NewRouter()` with `httptest`, so they need no running server.
- Middleware tests redirect `log` output to a buffer and restore it with
  `defer log.SetOutput(os.Stderr)`.

## Docker deployment

- **nginx as reverse proxy, same origin, no CORS.** The frontend container
  serves the built SPA and proxies `/api/` to `backend:8080`. The browser only
  talks to one origin, so CORS never comes into play in Docker. The backend
  publishes no ports to the host (only `frontend` maps `3000:80`), so it is
  reachable only inside the Compose network.
- **Multi-stage builds.** Both images build in a throwaway stage (Go / Node)
  and copy only the artifact into a slim runtime image (Alpine for the backend,
  nginx for the frontend).
- **Static Go binary, non-root user.** `CGO_ENABLED=0` with `-trimpath` and
  stripped symbols gives a self-contained binary; the backend image runs as a
  dedicated unprivileged user (`app`, uid 10001). The frontend image keeps the
  stock nginx user setup (master process as root) and is not hardened further.
- **Healthcheck controls start order.** The backend defines a healthcheck
  (`wget /health`) and the frontend uses `depends_on: condition: service_healthy`,
  so nginx only starts once the API answers. This covers startup only; nginx
  resolves `backend` at startup, so it must be restarted if the backend
  container is recreated with a new address.
- **`ENTRYPOINT` in exec form** (`["/server"]`). The server is PID 1 and
  receives `SIGTERM` directly from `docker compose down`, so the graceful
  shutdown in `main.go` runs. The shell form would wrap it in `sh -c`, which
  does not forward signals.
