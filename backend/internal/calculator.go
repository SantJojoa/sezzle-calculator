package calculator

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero = errors.New("cannot divide by zero")
	ErrNegativeNumber = errors.New("cannot calculate square root of negative number")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInvalidResult  = errors.New("invalid result")
)

func Add(a, b float64) (float64, error) {
	return finite(a + b)
}

func Subtract(a, b float64) (float64, error) {
	return finite(a - b)
}

func Multiply(a, b float64) (float64, error) {
	return finite(a * b)
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return finite(a / b)
}

func Power(base, exp float64) (float64, error) {
	return finite(math.Pow(base, exp))
}

func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeNumber
	}
	return finite(math.Sqrt(a))
}

func Percentage(value, pct float64) (float64, error) {
	return finite(value * pct / 100)
}

func finite(r float64) (float64, error) {
	if math.IsNaN(r) || math.IsInf(r, 0) {
		return 0, ErrInvalidResult
	}
	return r, nil
}
