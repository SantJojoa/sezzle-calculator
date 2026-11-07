// Package api exposes the calculator over HTTP as a JSON REST API.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/santiagojojoa/sizzle-calculator/backend/internal/calculator"
)

const (
	maxBodyBytes     = 4 << 10 // 4 KB: far more than any valid request needs
	maxExpressionLen = 200
)

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
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorBody{Code: code, Message: message}})
}

// decode reads exactly one JSON object into dst. It rejects oversized
// bodies, unknown fields, and trailing data like `{"a":1}{"b":2}`.
func decode(w http.ResponseWriter, r *http.Request, dst any) error {
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

func writeCalError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", err.Error())

	case errors.Is(err, calculator.ErrNegativeNumber):
		writeError(w, http.StatusUnprocessableEntity, "NEGATIVE_SQRT", err.Error())

	case errors.Is(err, calculator.ErrInvalidResult):
		writeError(w, http.StatusUnprocessableEntity, "INVALID_RESULT", err.Error())

	case errors.Is(err, calculator.ErrInvalidExpression):
		writeError(w, http.StatusUnprocessableEntity, "INVALID_EXPRESSION", err.Error())

	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
	}
}

type binaryOp func(float64, float64) (float64, error)

func handleBinary(op binaryOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req binaryRequest
		if err := decode(w, r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
			return
		}

		if req.A == nil || req.B == nil {
			writeError(w, http.StatusBadRequest, "MISSING_FIELDS", "Both 'a' and 'b' are required")
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
		writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
		return
	}

	if req.A == nil {
		writeError(w, http.StatusBadRequest, "MISSING_FIELDS", "'a' is required")
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
		writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
		return
	}

	if req.Expression == "" {
		writeError(w, http.StatusBadRequest, "MISSING_FIELDS", "'expression' is required")
		return
	}
	if len(req.Expression) > maxExpressionLen {
		writeError(w, http.StatusBadRequest, "EXPRESSION_TOO_LONG", "Expression exceeds the maximum length")
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
