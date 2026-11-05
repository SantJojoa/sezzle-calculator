package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	test := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 2, 3, 5, nil},
		{"negative", -2, -3, -5, nil},
		{"mixed-signs", -2, 5, 3, nil},
		{"with-zero", 7, 0, 7, nil},
		{"decimal", 1.5, 2.25, 3.75, nil},
		{"float-precision", 0.1, 0.2, 0.30000000000000004, nil},
		{"large-value-no-intermediate-overflow", math.MaxFloat64, 50, math.MaxFloat64 / 2, nil},
		{"max-float-no-overflow", math.MaxFloat64, 0, math.MaxFloat64, nil},
		{"overflow", math.MaxFloat64, math.MaxFloat64, 0, ErrInvalidResult},
		{"negative-overflow", -math.MaxFloat64, -math.MaxFloat64, 0, ErrInvalidResult},
		{"infinity-operand", math.Inf(1), 1, 0, ErrInvalidResult},
		{"infinity-minus-infinity", math.Inf(1), math.Inf(-1), 0, ErrInvalidResult},
		{"nan-operand", math.NaN(), 1, 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	test := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 10, 4, 6, nil},
		{"negative-result", 4, 10, -6, nil},
		{"negative-operands", -4, -10, 6, nil},
		{"same-operands", 5, 5, 0, nil},
		{"with-zero", 7, 0, 7, nil},
		{"decimal", 5.5, 2.25, 3.25, nil},
		{"float-precision", 0.3, 0.1, 0.19999999999999998, nil},
		{"overflow", math.MaxFloat64, -math.MaxFloat64, 0, ErrInvalidResult},
		{"negative-overflow", -math.MaxFloat64, math.MaxFloat64, 0, ErrInvalidResult},
		{"infinity-minus-infinity", math.Inf(1), math.Inf(1), 0, ErrInvalidResult},
		{"nan-operand", 1, math.NaN(), 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subtract(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	test := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 3, 4, 12, nil},
		{"by-zero", 5, 0, 0, nil},
		{"negative-times-positive", -3, 4, -12, nil},
		{"negative-times-negative", -3, -4, 12, nil},
		{"by-one", 9, 1, 9, nil},
		{"decimal", 0.5, 0.5, 0.25, nil},
		{"negative-zero", math.Copysign(0, -1), 5, 0, nil},
		{"overflow", math.MaxFloat64, 2, 0, ErrInvalidResult},
		{"negative-overflow", math.MaxFloat64, -2, 0, ErrInvalidResult},
		{"infinity-times-zero", math.Inf(1), 0, 0, ErrInvalidResult},
		{"nan-operand", math.NaN(), 2, 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Multiply(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	test := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 10, 2, 5, nil},
		{"by-zero", 10, 0, 0, ErrDivisionByZero},
		{"by-negative-zero", 10, math.Copysign(0, -1), 0, ErrDivisionByZero},
		{"zero-by-zero", 0, 0, 0, ErrDivisionByZero},
		{"zero-numerator", 0, 5, 0, nil},
		{"negative-numerator", -10, 4, -2.5, nil},
		{"negative-denominator", 10, -4, -2.5, nil},
		{"negative-operands", -10, -4, 2.5, nil},
		{"decimal", 1, 4, 0.25, nil},
		{"overflow", math.MaxFloat64, 0.5, 0, ErrInvalidResult},
		{"tiny-denominator-overflow", 1, math.SmallestNonzeroFloat64, 0, ErrInvalidResult},
		{"infinity-by-infinity", math.Inf(1), math.Inf(1), 0, ErrInvalidResult},
		{"nan-denominator", 1, math.NaN(), 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	test := []struct {
		name      string
		base, exp float64
		want      float64
		wantErr   error
	}{
		{"positive", 2, 10, 1024, nil},
		{"zero-exponent", 5, 0, 1, nil},
		{"zero-to-zero", 0, 0, 1, nil},
		{"zero-base", 0, 3, 0, nil},
		{"exponent-one", 7, 1, 7, nil},
		{"negative-exponent", 2, -1, 0.5, nil},
		{"negative-base-even-exponent", -2, 2, 4, nil},
		{"negative-base-odd-exponent", -2, 3, -8, nil},
		{"fractional-exponent", 9, 0.5, 3, nil},
		{"negative-base-fractional-exponent", -8, 0.5, 0, ErrInvalidResult},
		{"zero-to-negative-exponent", 0, -1, 0, ErrInvalidResult},
		{"overflow", 2, 1024, 0, ErrInvalidResult},
		{"negative-overflow", -2, 1025, 0, ErrInvalidResult},
		{"infinity-base", math.Inf(1), 2, 0, ErrInvalidResult},
		{"nan-base", math.NaN(), 2, 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Power(tt.base, tt.exp)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	test := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{"positive", 9, 3, nil},
		{"zero", 0, 0, nil},
		{"negative-zero", math.Copysign(0, -1), 0, nil},
		{"one", 1, 1, nil},
		{"decimal", 0.25, 0.5, nil},
		{"irrational-result", 2, math.Sqrt2, nil},
		{"max-float", math.MaxFloat64, 1.3407807929942596e154, nil},
		{"negative", -4, 0, ErrNegativeNumber},
		{"negative-infinity", math.Inf(-1), 0, ErrNegativeNumber},
		{"positive-infinity", math.Inf(1), 0, ErrInvalidResult},
		{"nan", math.NaN(), 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sqrt(tt.a)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	test := []struct {
		name       string
		value, pct float64
		want       float64
		wantErr    error
	}{
		{"positive", 200, 50, 100, nil},
		{"hundred-percent", 80, 100, 80, nil},
		{"zero-percent", 80, 0, 0, nil},
		{"zero-value", 0, 25, 0, nil},
		{"over-hundred-percent", 50, 200, 100, nil},
		{"decimal-percent", 80, 12.5, 10, nil},
		{"decimal-value", 12.5, 10, 1.25, nil},
		{"negative-value", -50, 10, -5, nil},
		{"negative-percent", 50, -10, -5, nil},
		{"negative-operands", -50, -10, 5, nil},
		{"overflow", math.MaxFloat64, math.MaxFloat64, 0, ErrInvalidResult},
		{"infinity-value", math.Inf(1), 10, 0, ErrInvalidResult},
		{"nan-percent", 50, math.NaN(), 0, ErrInvalidResult},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Percentage(tt.value, tt.pct)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
