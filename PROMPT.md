# AI Usage Log

This document lists the AI prompts used during development, as requested in the
assignment. Each entry says what the prompt was for, who wrote what, and what I
accepted, changed or rejected.

**Tools used:**
- Claude (Anthropic): chat for planning and design discussion, Claude Code for
  implementation help and reviews.
- Stitch (Google): generated the UI mockup used as the visual reference
  (section 8).

## How I used AI

I was coming from a TypeScript/NestJS background, so I used AI in three roles:

- A **mentor** for planning, Go conventions and design trade-offs.
- A **reviewer** to find edge cases and issues in code I had already written.
- A **generator** for repetitive test cases and for filling in pieces I had
  specified first (skeletons, signatures, test tables, `TODO` lists).

My usual pattern was to write the structure and the first examples myself, then
ask AI to follow that pattern. Each entry below lists exactly which parts I
wrote and which the AI generated. Every AI suggestion was reviewed, run and
adjusted before committing.

| # | Topic | Main use |
|---|-------|----------|
| 1 | Calculator package edge cases and tests | Reviewer, generator |
| 2 | Expression evaluator (parser) | Mentor, generator |
| 3 | Security and code review of the HTTP layer | Reviewer |
| 4 | HTTP router and middleware tests | Generator |
| 5 | Frontend API client and its tests | Generator |
| 6 | Frontend expression logic (`src/lib`) | Generator |
| 7 | React hook and components | Generator |
| 8 | Visual design (Stitch) and styling | Generator |
| 9 | Accessibility review of the React components | Reviewer |
| 10 | Component tests with React Testing Library | Generator |
| 11 | Closing coverage gaps | Generator |
| 12 | Docker setup review | Reviewer |
| 13 | Documentation, screenshots and coverage reports | Generator, reviewer |

---

## 1. Edge cases and unit tests for the calculator package

**Context:** I wrote `calculator.go` and the table-driven tests for `Add` and
`Divide` first.

**Prompt:**

> I'm sharing my Go calculator package and two existing table-driven unit tests
> for the Add and Divide functions.
>
> Act as a senior Go software engineer and perform a thorough review of my
> implementation and test coverage.
>
> Please:
>
> 1. Identify missing edge cases for every calculator operation, including:
>    - Integer and floating-point overflow.
>    - Division by zero.
>    - NaN and positive/negative infinity.
>    - Negative zero and floating-point precision.
>    - Negative operands and decimal values.
>    - Invalid inputs and boundary conditions.
>    - Square root of negative numbers and other operation-specific edge cases.
>
> 2. Explain which cases are relevant to my implementation and which are
>    unnecessary or redundant. Do not introduce tests for behavior my API does
>    not support.
>
> 3. Write table-driven unit tests for all remaining functions, following
>    exactly the naming conventions, structure, assertion style, and conventions
>    of my existing Add and Divide tests.
>
> 4. Use Go's standard testing package. Avoid unnecessary third-party
>    dependencies.
>
> 5. Do not modify my production code unless you identify a concrete bug. If a
>    change is necessary, explain why and suggest the smallest possible fix.
>
> 6. Ensure the tests are deterministic, independent, readable, and maintainable.
>
> 7. Provide the complete test code and briefly explain the purpose of each new
>    test case.
>
> Prioritize correctness, idiomatic Go practices, and meaningful test coverage
> over maximizing the number of test cases.

**Accepted:**
- Edge cases for overflow, NaN/±Inf, negative zero, and float precision
  (e.g. `0.1 + 0.2 = 0.30000000000000004`).
- Test tables for Subtract, Multiply, Power, Sqrt, and Percentage.

**Changed / fixed after review:**
- Fixed a mismatched error name between tests and production code
  (`ErrNegativeNumber` vs `ErrNegativeSqrt`).
- Renamed the table variable `test` → `tests` (idiomatic Go).
- Moved files into `internal/calculator/` so the package name matches its
  directory.
- Fixed a potential intermediate overflow in `Percentage` found in the AI
  review: changed `value * pct / 100` to `value * (pct / 100)`. I added a test
  case and confirmed it failed with the old formula before applying the fix.

**Note:** NaN/Inf cases can't arrive through the JSON API (JSON can't encode
them), but they protect the package's contract on its own and cover overflow
produced internally (e.g. `2^1024`).

**Result:** 100% statement coverage on `internal/calculator` at this point
(before the parser in section 2 was added).

---

## 2. Expression evaluator (parser)

**Purpose:** Support full expressions (`2+3*4`, parentheses, unary minus)
instead of only two-operand operations.

**What AI provided (in a guided conversation):**
- The design options (frontend chaining vs. backend parser) and the grammar.
- The test table for `Evaluate`, which I typed in and ran first (test-first).
- A skeleton with the token types, `Evaluate`, parser helpers, and
  `expression()` as a pattern to follow.
- `unary()` and `power()`, with an explanation of how they produce correct
  precedence and right-associativity.

**What I wrote:**
- The tokenizer logic (operators, number scanning, identifier scanning).
- `term()` (following the `expression()` pattern) and `primary()`.

**Result:** All 34 expression cases passed. Coverage on the package: 93%.

---

## 3. Security and code review of the HTTP layer (Claude Code)

**Prompt:** "Review this Go HTTP handler code for idiomatic error handling,
input validation and security issues. Don't rewrite it — list concrete problems
with a short explanation each."

**Accepted and fixed:**
- Consistent status-code rule: 400 = malformed input (incl. expression syntax),
  422 = valid input but mathematically undefined.
- Stopped returning raw decode errors to clients; 413 for oversized bodies.
- Consistent `MISSING_FIELD` code and message format; whitespace-only
  expressions rejected.
- Logged unexpected errors; removed unused `ErrInvalidInput`.
- Log injection (`%q`), `nosniff`, `Vary: Origin`, `IdleTimeout`.
- Status code in request logs, panic-recovery middleware, graceful shutdown,
  and `Content-Type: application/json` enforcement (415).

**Documented as known limitations (not fixed, with reasons):**
- Recursion depth bounded by the 200-byte cap rather than an explicit counter.
- `len()` counts bytes (tokenizer is ASCII-only; message now says "200 bytes").
- Plain-text 404/405 from the standard mux.
- CORS preflight answers on any path.
- `writeJSON` without a buffer; handler duplication; no `slog`/DI;
  `ALLOWED_ORIGIN="*"` not validated.
  (Reasons for each are in the README.)

**Note:** Some findings corrected suggestions from an earlier AI-guided
session (returning decode errors verbatim, 422 for syntax errors).

---

## 4. HTTP router and middleware tests (Claude Code)

**Prompt:** "Here's my Go HTTP router test with a table-driven structure and
two cases I wrote. Add cases following exactly the same style for: each
endpoint's success path; a zero operand to verify missing-field detection uses
pointers; 400 cases (invalid JSON, empty body, unknown field, trailing data,
missing a/b, string instead of number, sqrt without a, empty and whitespace-only
expression, expression over 200 bytes, invalid syntax); 413 for a body over
4 KB; 422 cases (negative sqrt, overflow, division by zero inside an
expression); 405 for GET on a POST route; 404 for an unknown path. Use my error
codes. Don't change the test loop."

**What I wrote:**
- The test loop and the first two router cases (`add-success`,
  `divide-by-zero`).
- The structure of the middleware tests (`TestCORSPreflight` skeleton).

**What AI generated:**
- The remaining 29 router cases, in the same positional-literal style.
- `TestCORSRegularRequest`, `TestRecover` and `TestLoggerRecordsStatus`, and
  the missing assertions in `TestCORSPreflight` (the file did not compile:
  `called` was declared but unused).

**What I adjusted:**
- Oversized-body cases are built from `maxExpressionLen` / `maxBodyBytes`
  instead of magic numbers; the 413 payload goes inside `expression` so the
  decoder hits the size limit before an unknown-field error.
- Log redirection in the middleware tests restores `os.Stderr` with `defer`
  rather than silencing the logger after the test.
- `cmd/server` is left untested on purpose (documented in the README) instead
  of adding an artificial test to raise the number.

**Result:** `internal/api` 96.3%, `internal/calculator` 93.0%.

---

## 5. Frontend API client and its tests (Claude Code)

**Prompt:** "Write Vitest tests for this API client, mocking global fetch with
vi.stubGlobal: successful result, backend JSON error (verify ApiError code,
message and status), non-JSON error response (text/plain 405), and network
failure. Restore mocks after each test."

**What I wrote:** the client skeleton in `calculatorClient.ts` (`ApiError`,
the `evaluate` request and network-error handling) and the `TODO` list for the
error and result branches.

**What AI generated:** the 4 tests in `calculatorClient.test.ts`, plus the
`TODO` bodies (error/result parsing with type guards) and the endpoint fix
(`/api/v1/calculate` → `/api/v1/evaluate`).

**Result:** 4/4 tests pass.

---

## 6. Frontend expression logic (`src/lib`) (Claude Code)

**Prompt:** "Write the tests first: `src/lib/expression.test.ts`", followed by
my table of cases for `appendKey`, `backspace`, `validate`, `toDisplay` and
`formatResult`. Then: "Look at my `expression.ts` and implement `validate`
(parenthesis counter), the dot rule in `appendKey` and `formatResult`, and
document the formatting decision."

**What I wrote:** the table of cases and expected values, the function
signatures with `TODO` markers, `backspace` and `toDisplay`.

**What AI generated:** the test file from my table, and the implementations of
`validate`, `appendKey` and the doc comment for `formatResult`.

**What I adjusted:** `appendKey` needed more rules than the dot rule to pass
the whole table (operator replacement, `-` after an operator, no implicit
multiplication, `)` only with an open `(`). I reviewed these against the
backend grammar, which rejects unary `+` and `2(3)`.

**Result:** 18 test rows pass.

---

## 7. React hook and components (Claude Code)

**Prompts:** one per file: "make the `useCalculator.ts`", then specs for
`Display.tsx`, `Keypad.tsx` and `Calculator.tsx` (props, key layout, which key
sends which API character, `aria-live` / `role="alert"`, physical keyboard
support, cleanup of the listener).

**What I wrote:** the hook outline (state shape, `press`, `clear`,
`deleteLast`, `submit` steps as comments) and the component specs.

**What AI generated:** the implementations of the hook and the three
components, and `App.tsx`.

**What I adjusted or decided:**
- `=` spans two rows in the keypad so the grid has no empty cell.
- A request counter in the hook ignores a slow response for an expression the
  user has already changed.
- After a result, an operator only continues from it if the backend can parse
  it back (exponent notation like `1e+21` cannot).

---

## 8. Visual design and styling (Claude Code)

**Prompts:**
1. A design prompt for the calculator UI (palette, layout, states), given to
   **Stitch (Google)** to generate the reference image. This one was not
   written with Claude.
2. A styling prompt for Claude Code with strict rules: apply visuals only, do not change
   structure, props, logic, `aria-*` attributes or tests, and use only the
   given theme tokens.

> Design a clean, modern calculator web app UI, single screen, shown on both
> mobile (375px wide) and desktop. Dark theme only. Page background `#211733`;
> operator keys `#693a94`; "=" key `#9ff4d9`; WCAG AA contrast; centered card
> (max width ~360px); 4-column keypad with `C ( ) ⌫ / √ ^ ÷ × / 7 8 9 − /
> 4 5 6 + / 1 2 3 = / 0 .`; states: initial, typing, result, error, loading;
> visible focus rings; no logos, no features beyond the calculator.

**What I wrote:** the design brief, the color tokens and the list of things to
ignore from the reference image (tabs, profile icon, bottom navigation, `%`
key, and so on).

**What AI generated:** the mockup (Stitch), and then the Tailwind classes and
`index.css` tokens that implement it (Claude Code).

**What I checked:** the rendered result at 375px against the reference, and
that the existing tests still passed.

---

## 9. Accessibility review of the React components (Claude Code)

**Prompt:** "Review this React component for accessibility (labels,
aria-live, keyboard usage, focus) and suggest minimal changes only."

**Accepted and applied (1-6 and 8):**
- Skip `Enter` when the event target is a button, so a focused key activates
  itself instead of submitting.
- `<title>Calculator</title>`.
- `role="status"` + `aria-atomic` on the always-mounted result region, with a
  visually hidden "Result: " prefix; `aria-busy` moved onto that region.
- `role="group"` + `aria-label` on the keypad; hidden "Expression: " prefix.
- Ignore shortcuts when focus is in an input or textarea.

**Left out (7), with the reason:** replacing `disabled` with `aria-disabled`
on `=`. The impact is low (the listener is on `window`, so Enter keeps working
when focus is lost) and it would force changing the "disabled while loading"
behaviour and its test. Documented under "Known limitations" in the README.

**Side effect of #1, found and fixed:** once Enter is ignored on buttons, a
mouse click leaves focus on the clicked key, so a following Enter re-pressed
that key instead of submitting. Fixed with `onMouseDown={e => e.preventDefault()}`
on the keypad buttons so clicks do not move focus.

**Tests added:** `Calculator.test.tsx` covers the live region (`status`),
Enter on a focused key, Enter after a mouse click, the group/symbol labels and
the error `alert`. I checked that the two Enter tests fail when their fix is
removed.

---

## 10. Component tests with React Testing Library (Claude Code)

**Prompt:** "Write React Testing Library tests for this Calculator component
using userEvent. Mock the `../api/calculatorClient` module with `vi.mock`.
Cover: clicking keys builds the expression; `=` shows the formatted result;
backend ApiError message is shown in the alert; client-side validation error
(unclosed parenthesis) does NOT call the API; pressing an operator after a
result continues from it; physical keyboard input (typing `2+3` and Enter);
`=` is disabled while loading."

**What I wrote:** the list of behaviors to cover.

**What AI generated:** the 8 tests added to `Calculator.test.tsx`.

**What I checked:** that the API is called with the API character (`2*3`), not
the display symbol, and that `0.1 + 0.2` is shown as `0.3`.

---

## 11. Closing coverage gaps (Claude Code)

**Prompt:** "Here are my React tests and the uncovered lines from the coverage
report: [`pnpm coverage` output]. Write the minimal tests to cover those
branches, following my existing style and using `getByRole` queries."

**What I wrote:** the `pnpm coverage` run and the decision to cover only the
reported lines.

**What AI generated:** tests for Backspace, Escape and the `C` / delete keys,
an `App` render test, three `appendKey` rows (`)` and unknown keys) and three
client tests for malformed responses.

**Result:** 46/46 frontend tests pass; statements 95.8%, branches 86.1%,
lines 100%. Lint is clean.

---

## 12. Docker setup review (Claude Code)

**Prompt:** "Check that everything is ready to deliver", which included
building and running the full stack.

**What AI did:** built both images with `docker compose up --build`, called the
API through nginx on port 3000, checked that the backend port is not published,
that the backend runs as a non-root user and that `SIGTERM` triggers the
graceful shutdown. It found that `GET /health` through nginx returned the app's
`index.html` instead of the backend response.

**What I decided:** to proxy `/health` in `nginx.conf` (a `location = /health`
block) rather than document it as local-only. Re-verified after rebuilding:
`/health` now returns `{"status":"ok"}` through port 3000.

---

## 13. Documentation, screenshots and coverage reports (Claude Code)

**Prompts:** "Improve clarity and structure of this README. Don't add any
claims that aren't in my draft.", "Now review PROMPT.md" and a final
delivery check against the assignment text.

**What I wrote:** the README structure (sections, tables, the 400/422 rule),
the decision of which design points to document, and the requirement that every
README claim be verifiable.

**What AI generated:** the filled-in README sections (from my notes and the
code), the verified coverage numbers and curl responses, the restructuring, the
review of this file, the committed coverage reports in `docs/coverage/`, and the
desktop and mobile screenshots (a script drives a headless browser against the
running Docker stack, types `2+3*4` and captures the result).

**What I checked:** that the README's numbers, endpoints and error codes match
the code and a real run.
