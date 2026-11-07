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
