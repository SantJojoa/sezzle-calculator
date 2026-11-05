package calculator

import (
	"fmt"
	"strconv"
)

type tokenKind int

const (
	tokNumber tokenKind = iota
	tokPlus
	tokMinus
	tokStar
	tokSlash
	tokCaret
	tokLParen
	tokRParen
	tokSqrt
	tokEOF
)

type token struct {
	kind  tokenKind
	value float64
	pos   int
}

func tokenize(input string) ([]token, error) {
	var tokens []token
	i := 0
	for i < len(input) {
		c := input[i]
		switch {
		case c == ' ' || c == '\t':
			i++
		case c == '+':
			tokens = append(tokens, token{kind: tokPlus, pos: i})
			i++
		case c == '-':
			tokens = append(tokens, token{kind: tokMinus, pos: i})
			i++
		case c == '*':
			tokens = append(tokens, token{kind: tokStar, pos: i})
			i++
		case c == '/':
			tokens = append(tokens, token{kind: tokSlash, pos: i})
			i++
		case c == '^':
			tokens = append(tokens, token{kind: tokCaret, pos: i})
			i++
		case c == '(':
			tokens = append(tokens, token{kind: tokLParen, pos: i})
			i++
		case c == ')':
			tokens = append(tokens, token{kind: tokRParen, pos: i})
			i++
		case isDigit(c) || c == '.':
			start := i
			for i < len(input) && (isDigit(input[i]) || input[i] == '.') {
				i++
			}
			n, err := strconv.ParseFloat(input[start:i], 64)
			if err != nil {
				return nil, invalidAt(start)
			}
			tokens = append(tokens, token{kind: tokNumber, value: n, pos: start})
		case isLetter(c):
			start := i
			for i < len(input) && isLetter(input[i]) {
				i++
			}
			if input[start:i] != "sqrt" {
				return nil, invalidAt(start)
			}
			tokens = append(tokens, token{kind: tokSqrt, pos: start})
		default:
			return nil, invalidAt(i)

		}
	}
	tokens = append(tokens, token{kind: tokEOF, pos: len(input)})
	return tokens, nil
}

type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token { return p.tokens[p.pos] }

func (p *parser) next() token {
	t := p.tokens[p.pos]
	if t.kind != tokEOF {
		p.pos++
	}
	return t
}

// expect consumes the next token if it has the given kind, or fails.
func (p *parser) expect(kind tokenKind) error {
	if t := p.peek(); t.kind != kind {
		return invalidAt(t.pos)
	}
	p.next()
	return nil
}

// Evaluate parses and computes an arithmetic expression such as "2+3*4".
// Supports + - * / ^, parentheses, unary minus and sqrt(...).
func Evaluate(input string) (float64, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return 0, err
	}
	p := &parser{tokens: tokens}
	result, err := p.expression()
	if err != nil {
		return 0, err
	}
	if t := p.peek(); t.kind != tokEOF { // leftover input, e.g. "2+3)" or "2(3)"
		return 0, invalidAt(t.pos)
	}
	return result, nil
}

// expression = term { ("+" | "-") term }
func (p *parser) expression() (float64, error) {
	left, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		switch p.peek().kind {
		case tokPlus:
			p.next()
			right, err := p.term()
			if err != nil {
				return 0, err
			}
			if left, err = Add(left, right); err != nil {
				return 0, err
			}
		case tokMinus:
			p.next()
			right, err := p.term()
			if err != nil {
				return 0, err
			}
			if left, err = Subtract(left, right); err != nil {
				return 0, err
			}
		default:
			return left, nil
		}
	}
}

// term = unary { ("*" | "/") unary }
func (p *parser) term() (float64, error) {
	left, err := p.unary()
	if err != nil {
		return 0, err
	}
	for {
		switch p.peek().kind {
		case tokStar:
			p.next()
			right, err := p.unary()
			if err != nil {
				return 0, err
			}
			if left, err = Multiply(left, right); err != nil {
				return 0, err
			}
		case tokSlash:
			p.next()
			right, err := p.unary()
			if err != nil {
				return 0, err
			}
			if left, err = Divide(left, right); err != nil {
				return 0, err
			}
		default:
			return left, nil
		}
	}
}

// unary = "-" unary | power
func (p *parser) unary() (float64, error) {
	if p.peek().kind == tokMinus {
		p.next()
		v, err := p.unary() // recursive: allows "--2"
		if err != nil {
			return 0, err
		}
		return -v, nil
	}
	return p.power()
}

// power = primary [ "^" unary ]
func (p *parser) power() (float64, error) {
	base, err := p.primary()
	if err != nil {
		return 0, err
	}
	if p.peek().kind == tokCaret {
		p.next()
		exp, err := p.unary() // unary, not primary: enables 2^-1 and right-associativity
		if err != nil {
			return 0, err
		}
		return Power(base, exp)
	}
	return base, nil
}

// primary = number | "(" expression ")" | "sqrt" "(" expression ")"
func (p *parser) primary() (float64, error) {
	t := p.next()
	switch t.kind {
	case tokNumber:
		return t.value, nil
	case tokLParen:
		v, err := p.expression()
		if err != nil {
			return 0, err
		}
		if err := p.expect(tokRParen); err != nil {
			return 0, err
		}
		return v, nil
	case tokSqrt:
		if err := p.expect(tokLParen); err != nil {
			return 0, err
		}
		v, err := p.expression()
		if err != nil {
			return 0, err
		}
		if err := p.expect(tokRParen); err != nil {
			return 0, err
		}
		return Sqrt(v)
	default:
		return 0, invalidAt(t.pos)
	}
}

func isDigit(c byte) bool  { return c >= '0' && c <= '9' }
func isLetter(c byte) bool { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }

func invalidAt(pos int) error {
	return fmt.Errorf("%w at position %d", ErrInvalidExpression, pos)
}
