package calculator

import (
	"errors"
	"testing"
)

func TestEvaluate(t *testing.T) {
	test := []struct {
		name       string
		expression string
		expected   float64
		err        error
	}{
		{"single-number", "42", 42, nil},
		{"addition", "2+3", 5, nil},
		{"precedence", "2+3*4", 14, nil},
		{"left-assoc-subtract", "10-4-3", 3, nil},
		{"left-assoc-divide", "100/10/2", 5, nil},
		{"parentheses", "(2+3)*4", 20, nil},
		{"nested-parentheses", "((1+2)*(3+4))", 21, nil},
		{"whitespace", "  2 +   3 ", 5, nil},
		{"decimals", "3.5*2", 7, nil},
		{"leading-dot", ".5+.5", 1, nil},
		// unary minus & power
		{"unary-minus", "-5+3", -2, nil},
		{"double-negative", "--2", 2, nil},
		{"negative-operand", "2*-3", -6, nil},
		{"power", "2^10", 1024, nil},
		{"power-right-assoc", "2^3^2", 512, nil},
		{"unary-minus-before-power", "-2^2", -4, nil},
		{"negative-exponent", "2^-1", 0.5, nil},
		// sqrt
		{"sqrt", "sqrt(16)", 4, nil},
		{"sqrt-expression", "sqrt(9+16)", 5, nil},
		// domain errors (reuse existing ones)
		{"division-by-zero", "10/(5-5)", 0, ErrDivisionByZero},
		{"sqrt-negative", "sqrt(-4)", 0, ErrNegativeNumber},
		{"overflow", "2^1024", 0, ErrInvalidResult},
		// syntax errors
		{"empty", "", 0, ErrInvalidExpression},
		{"only-spaces", "   ", 0, ErrInvalidExpression},
		{"double-operator", "2+*3", 0, ErrInvalidExpression},
		{"trailing-operator", "2+", 0, ErrInvalidExpression},
		{"unclosed-paren", "(2+3", 0, ErrInvalidExpression},
		{"extra-close-paren", "2+3)", 0, ErrInvalidExpression},
		{"empty-parens", "()", 0, ErrInvalidExpression},
		{"implicit-multiplication", "2(3)", 0, ErrInvalidExpression},
		{"malformed-number", "1.2.3", 0, ErrInvalidExpression},
		{"unknown-function", "abc(2)", 0, ErrInvalidExpression},
		{"sqrt-without-parens", "sqrt16", 0, ErrInvalidExpression},
		{"invalid-character", "2&3", 0, ErrInvalidExpression},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression)
			if !errors.Is(err, tt.err) {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
			if got != tt.expected {
				t.Fatalf("got = %v, want %v", got, tt.expected)
			}
		})
	}
}
