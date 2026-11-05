# AI Usage Log

This document lists the AI prompts used during development, as requested in the
assignment. For each one I note what it was used for and what I accepted,
changed, or rejected.

**Tools used:** Claude (Anthropic)

## How I used AI

I wrote the core production code (calculator logic, HTTP handlers, expression
parser, React components) myself. I used AI as:

- A **mentor** for planning, Go conventions (I'm coming from a TypeScript/NestJS
  background), and design trade-offs.
- A **reviewer** to find edge cases and issues in code I had already written.
- A **generator** for repetitive test cases, following patterns I wrote first.

Every AI-generated suggestion was reviewed, run, and adjusted before committing.


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
- [Fixed a mismatched error name between tests and production code
  (`ErrNegativeNumber` vs `ErrNegativeSqrt`).]
- [Renamed the table variable `test` → `tests` (idiomatic Go).]
- [Moved files into `internal/calculator/` so the package name matches its directory.]
- [Fixed a potential intermediate overflow in `Percentage`: changed
  `value * pct / 100` to `value * (pct / 100)` and added a test case for it.]

**Note:** NaN/Inf cases can't arrive through the JSON API (JSON can't encode
them), but they protect the package's contract on its own and cover overflow
produced internally (e.g. `2^1024`).

**Result:** 100% statement coverage on `internal/calculator`.

---

<!-- Add one section per prompt as you go: handler review, handler tests,
     parser help, frontend tests, README polish... -->