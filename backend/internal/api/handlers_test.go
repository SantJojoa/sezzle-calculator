package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func ptr(v float64) *float64 { return &v }

func TestRouter(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantResult *float64 // set only for successful calculations
		wantCode   string   // set only for JSON error responses
	}{
		{"add-success", http.MethodPost, "/api/v1/add", `{"a":5,"b":3}`, http.StatusOK, ptr(8), ""},
		{"divide-by-zero", http.MethodPost, "/api/v1/divide", `{"a":5,"b":0}`, http.StatusUnprocessableEntity, nil, "DIVISION_BY_ZERO"},

		// success path for every endpoint
		{"subtract-success", http.MethodPost, "/api/v1/subtract", `{"a":5,"b":3}`, http.StatusOK, ptr(2), ""},
		{"multiply-success", http.MethodPost, "/api/v1/multiply", `{"a":4,"b":3}`, http.StatusOK, ptr(12), ""},
		{"divide-success", http.MethodPost, "/api/v1/divide", `{"a":10,"b":4}`, http.StatusOK, ptr(2.5), ""},
		{"power-success", http.MethodPost, "/api/v1/power", `{"a":2,"b":10}`, http.StatusOK, ptr(1024), ""},
		{"percentage-success", http.MethodPost, "/api/v1/percentage", `{"a":200,"b":15}`, http.StatusOK, ptr(30), ""},
		{"sqrt-success", http.MethodPost, "/api/v1/sqrt", `{"a":16}`, http.StatusOK, ptr(4), ""},
		{"evaluate-success", http.MethodPost, "/api/v1/evaluate", `{"expression":"2+3*4"}`, http.StatusOK, ptr(14), ""},
		{"health-success", http.MethodGet, "/health", ``, http.StatusOK, nil, ""},

		// zero operands are valid values, not missing fields (pointers)
		{"add-zero-operands", http.MethodPost, "/api/v1/add", `{"a":0,"b":0}`, http.StatusOK, ptr(0), ""},
		{"sqrt-zero-operand", http.MethodPost, "/api/v1/sqrt", `{"a":0}`, http.StatusOK, ptr(0), ""},

		// 400: malformed input
		{"invalid-json", http.MethodPost, "/api/v1/add", `{"a":`, http.StatusBadRequest, nil, "INVALID_JSON"},
		{"empty-body", http.MethodPost, "/api/v1/add", ``, http.StatusBadRequest, nil, "INVALID_JSON"},
		{"unknown-field", http.MethodPost, "/api/v1/add", `{"a":5,"b":3,"c":1}`, http.StatusBadRequest, nil, "INVALID_JSON"},
		{"trailing-data", http.MethodPost, "/api/v1/add", `{"a":5,"b":3}{"a":1,"b":1}`, http.StatusBadRequest, nil, "INVALID_JSON"},
		{"missing-a", http.MethodPost, "/api/v1/add", `{"b":3}`, http.StatusBadRequest, nil, "MISSING_FIELD"},
		{"missing-b", http.MethodPost, "/api/v1/add", `{"a":5}`, http.StatusBadRequest, nil, "MISSING_FIELD"},
		{"string-instead-of-number", http.MethodPost, "/api/v1/add", `{"a":"5","b":3}`, http.StatusBadRequest, nil, "INVALID_JSON"},
		{"sqrt-missing-a", http.MethodPost, "/api/v1/sqrt", `{}`, http.StatusBadRequest, nil, "MISSING_FIELD"},
		{"evaluate-empty-expression", http.MethodPost, "/api/v1/evaluate", `{"expression":""}`, http.StatusBadRequest, nil, "MISSING_FIELD"},
		{"evaluate-whitespace-expression", http.MethodPost, "/api/v1/evaluate", `{"expression":"   "}`, http.StatusBadRequest, nil, "MISSING_FIELD"},
		{"evaluate-expression-too-long", http.MethodPost, "/api/v1/evaluate", `{"expression":"` + strings.Repeat("1", maxExpressionLen+1) + `"}`, http.StatusBadRequest, nil, "EXPRESSION_TOO_LONG"},
		{"evaluate-invalid-syntax", http.MethodPost, "/api/v1/evaluate", `{"expression":"2+*3"}`, http.StatusBadRequest, nil, "INVALID_EXPRESSION"},

		// 413: body over 4 KB
		{"body-too-large", http.MethodPost, "/api/v1/evaluate", `{"expression":"` + strings.Repeat("1", maxBodyBytes+1) + `"}`, http.StatusRequestEntityTooLarge, nil, "BODY_TOO_LARGE"},

		// 422: valid input, undefined math
		{"sqrt-negative", http.MethodPost, "/api/v1/sqrt", `{"a":-4}`, http.StatusUnprocessableEntity, nil, "NEGATIVE_SQRT"},
		{"power-overflow", http.MethodPost, "/api/v1/power", `{"a":10,"b":1000}`, http.StatusUnprocessableEntity, nil, "INVALID_RESULT"},
		{"evaluate-divide-by-zero", http.MethodPost, "/api/v1/evaluate", `{"expression":"10/0"}`, http.StatusUnprocessableEntity, nil, "DIVISION_BY_ZERO"},

		// router-level responses (status only)
		{"wrong-method", http.MethodGet, "/api/v1/add", ``, http.StatusMethodNotAllowed, nil, ""},
		{"unknown-path", http.MethodPost, "/api/v1/unknown", `{"a":1,"b":2}`, http.StatusNotFound, nil, ""},
	}

	router := NewRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}

			switch {
			case tt.wantResult != nil:
				var got resultResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode result: %v", err)
				}
				if got.Result != *tt.wantResult {
					t.Errorf("result = %v, want %v", got.Result, *tt.wantResult)
				}
			case tt.wantCode != "":
				var got errorResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode error: %v", err)
				}
				if got.Error.Code != tt.wantCode {
					t.Errorf("code = %q, want %q", got.Error.Code, tt.wantCode)
				}
			}
			// cases with neither (405, 404) only check the status code
		})
	}
}
func TestRejectsNonJSONContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/add", strings.NewReader(`{"a":1,"b":2}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnsupportedMediaType)
	}
	var got errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Error.Code != "UNSUPPORTED_MEDIA_TYPE" {
		t.Errorf("code = %q, want %q", got.Error.Code, "UNSUPPORTED_MEDIA_TYPE")
	}
}
