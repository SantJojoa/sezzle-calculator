import { describe, expect, it } from 'vitest'
import { appendKey, backspace, formatResult, toDisplay, validate } from './expression'

describe('appendKey', () => {
    it.each([
        { name: 'appends a digit', expr: '12', key: '3', want: '123' },
        { name: 'replaces a trailing operator', expr: '2+', key: '*', want: '2*' },
        { name: 'allows - after an operator (negative number)', expr: '2*', key: '-', want: '2*-' },
        { name: 'rejects * at the start', expr: '', key: '*', want: '' },
        { name: 'rejects / at the start', expr: '', key: '/', want: '' },
        { name: 'rejects ^ at the start', expr: '', key: '^', want: '' },
        { name: 'rejects a second dot in the same number', expr: '3.1', key: '.', want: '3.1' },
        { name: 'starts a decimal with 0.', expr: '', key: '.', want: '0.' },
        { name: 'ignores input past 200 characters', expr: '1'.repeat(200), key: '1', want: '1'.repeat(200) },
    ])('$name: "$expr" + "$key" -> "$want"', ({ expr, key, want }) => {
        expect(appendKey(expr, key)).toBe(want)
    })
})

describe('backspace', () => {
    it.each([
        { name: 'removes the last character', expr: '123', want: '12' },
        { name: 'removes "sqrt(" as a whole', expr: '2+sqrt(', want: '2+' },
    ])('$name: "$expr" -> "$want"', ({ expr, want }) => {
        expect(backspace(expr)).toBe(want)
    })
})

describe('validate', () => {
    it.each([
        { name: 'empty expression', expr: '', want: 'Enter an expression' },
        { name: 'unclosed parenthesis', expr: '(2+3', want: 'Missing closing parenthesis' },
        { name: 'extra closing parenthesis', expr: '2+3)', want: 'Unexpected closing parenthesis' },
        { name: 'valid expression', expr: '2+3', want: null },
    ])('$name: "$expr"', ({ expr, want }) => {
        expect(validate(expr)).toBe(want)
    })
})

describe('toDisplay', () => {
    it('shows operators and sqrt with math symbols', () => {
        expect(toDisplay('2*3/sqrt(4)')).toBe('2×3÷√(4)')
    })
})

describe('formatResult', () => {
    it.each([
        { name: 'hides floating point noise', value: 0.1 + 0.2, want: '0.3' },
        { name: 'keeps exponent notation for huge numbers', value: 1e21, want: '1e+21' },
    ])('$name: $value -> "$want"', ({ value, want }) => {
        expect(formatResult(value)).toBe(want)
    })
})
