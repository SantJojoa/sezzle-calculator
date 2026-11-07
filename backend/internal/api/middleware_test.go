package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const testOrigin = "http://localhost:5173"

func TestCORSPreflight(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/add", nil)
	rec := httptest.NewRecorder()
	CORS(testOrigin, next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if called {
		t.Error("next handler was called on a preflight request")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != testOrigin {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, testOrigin)
	}
}

func TestCORSRegularRequest(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodPost, "/api/v1/add", nil)
	rec := httptest.NewRecorder()
	CORS(testOrigin, next).ServeHTTP(rec, req)

	if !called {
		t.Error("next handler was not called on a regular request")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != testOrigin {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, testOrigin)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Errorf("Vary = %q, want %q", got, "Origin")
	}
}

func TestRecover(t *testing.T) {
	// silence the expected "panic: boom" log line
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr) // restore when done

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("boom") })

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Recover(next).ServeHTTP(rec, req) // reaching the next line means the panic was recovered

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	var got errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error.Code != "INTERNAL_ERROR" {
		t.Errorf("code = %q, want %q", got.Error.Code, "INTERNAL_ERROR")
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("panic was not logged: %q", buf.String())
	}
}

func TestLoggerRecordsStatus(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr) // restore when done

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/brew", nil)
	rec := httptest.NewRecorder()
	Logger(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
	if !strings.Contains(buf.String(), "418") {
		t.Errorf("log %q does not contain status 418", buf.String())
	}
}
