// Package api exposes the calculator over HTTP as a JSON REST API.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/santiagojojoa/sizzle-calculator/backend/internal/calculator"
)

const (
	maxBodyBytes     = 4 << 10 // 4 KB: far more than any valid request needs
	maxExpressionLen = 200
)

var errUnsupportedMediaType = errors.New("content type must be application/json")

// --- Request / response types ---

// Pointers distinguish "field missing" (nil) from "field is 0".
type binaryRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type unaryRequest struct {
	A *float64 `json:"a"`
}

type expressionRequest struct {
	Expression string `json:"expression"`
}

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

// --- JSON helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorBody{Code: code, Message: message}})
}

// decode reads exactly one JSON object into dst. It rejects oversized
// bodies, unknown fields, and trailing data like `{"a":1}{"b":2}`.
func decode(w http.ResponseWriter, r *http.Request, dst any) error {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		return errUnsupportedMediaType
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

// writeDecodeError maps JSON decoding failures to safe client messages.
func writeDecodeError(w http.ResponseWriter, err error) {
	var maxErr *http.MaxBytesError
	switch {
	case errors.Is(err, errUnsupportedMediaType):
		writeError(w, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
	case errors.As(err, &maxErr):
		writeError(w, http.StatusRequestEntityTooLarge, "BODY_TOO_LARGE", "request body must not exceed 4 KB")
	case strings.HasPrefix(err.Error(), "json: unknown field"):
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body contains an unknown field")
	default:
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must be a single valid JSON object")
	}
}

// writeMissingField reports a required field that is absent or empty.
func writeMissingField(w http.ResponseWriter, field string) {
	writeError(w, http.StatusBadRequest, "MISSING_FIELD", `field "`+field+`" is required`)
}

func writeCalError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", err.Error())

	case errors.Is(err, calculator.ErrNegativeNumber):
		writeError(w, http.StatusUnprocessableEntity, "NEGATIVE_SQRT", err.Error())

	case errors.Is(err, calculator.ErrInvalidResult):
		writeError(w, http.StatusUnprocessableEntity, "INVALID_RESULT", err.Error())

	case errors.Is(err, calculator.ErrInvalidExpression):
		writeError(w, http.StatusBadRequest, "INVALID_EXPRESSION", err.Error())

	default:
		log.Printf("unexpected calculator error: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
	}
}

type binaryOp func(float64, float64) (float64, error)

func handleBinary(op binaryOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req binaryRequest
		if err := decode(w, r, &req); err != nil {
			writeDecodeError(w, err)
			return
		}

		if req.A == nil {
			writeMissingField(w, "a")
			return
		}
		if req.B == nil {
			writeMissingField(w, "b")
			return
		}

		result, err := op(*req.A, *req.B)
		if err != nil {
			writeCalError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, resultResponse{Result: result})
	}
}

func handleSqrt(w http.ResponseWriter, r *http.Request) {
	var req unaryRequest
	if err := decode(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if req.A == nil {
		writeMissingField(w, "a")
		return
	}

	result, err := calculator.Sqrt(*req.A)
	if err != nil {
		writeCalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: result})
}

func handleEvaluate(w http.ResponseWriter, r *http.Request) {
	var req expressionRequest
	if err := decode(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if strings.TrimSpace(req.Expression) == "" {
		writeMissingField(w, "expression")
		return
	}
	if len(req.Expression) > maxExpressionLen {
		writeError(w, http.StatusBadRequest, "EXPRESSION_TOO_LONG", "expression must not exceed 200 bytes")
		return
	}

	result, err := calculator.Evaluate(req.Expression)
	if err != nil {
		writeCalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: result})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// NewRouter registers all API routes.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /api/v1/add", handleBinary(calculator.Add))
	mux.HandleFunc("POST /api/v1/subtract", handleBinary(calculator.Subtract))
	mux.HandleFunc("POST /api/v1/multiply", handleBinary(calculator.Multiply))
	mux.HandleFunc("POST /api/v1/divide", handleBinary(calculator.Divide))
	mux.HandleFunc("POST /api/v1/power", handleBinary(calculator.Power))
	mux.HandleFunc("POST /api/v1/percentage", handleBinary(calculator.Percentage))
	mux.HandleFunc("POST /api/v1/sqrt", handleSqrt)
	mux.HandleFunc("POST /api/v1/evaluate", handleEvaluate)
	return mux
}
