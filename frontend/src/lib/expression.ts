export const MAX_LENGTH = 200
export const OPERATORS = ['+', '-', '*', '/', '^'] as const

const isOperator = (c: string) => (OPERATORS as readonly string[]).includes(c)
const isDigit = (c: string) => c >= '0' && c <= '9'

/** Appends a key to the expression applying basic input rules. */
export function appendKey(expr: string, key: string): string {
    if (expr.length + key.length > MAX_LENGTH) return expr

    const last = expr.slice(-1)

    if (key === '.') {
        // current number = digits/dots after the last operator or parenthesis
        const current = /[\d.]*$/.exec(expr)?.[0] ?? ''
        if (current.includes('.')) return expr
        if (last === ')') return expr
        return current === '' ? expr + '0.' : expr + '.'
    }

    if (isDigit(key)) {
        return last === ')' ? expr : expr + key // no implicit multiplication: 2(3)
    }

    if (isOperator(key)) {
        // a leading or post-"(" operator is only valid as a negative sign
        if (expr === '' || last === '(') return key === '-' ? expr + key : expr
        if (isOperator(last)) {
            if (key === '-') return last === '-' ? expr : expr + key // 2*-3
            // replace the trailing operator(s), unless nothing valid precedes them
            const base = expr.replace(/[+\-*/^]+$/, '')
            return base === '' ? expr : base + key
        }
        return expr + key
    }

    if (key === '(' || key === 'sqrt(') {
        return isDigit(last) || last === ')' || last === '.' ? expr : expr + key
    }

    if (key === ')') {
        const open = [...expr].reduce((n, c) => n + (c === '(' ? 1 : c === ')' ? -1 : 0), 0)
        return open > 0 && !isOperator(last) && last !== '(' ? expr + key : expr
    }

    return expr
}

/** Removes the last key; removes "sqrt(" as a single unit. */
export function backspace(expr: string): string {
    if (expr.endsWith('sqrt(')) return expr.slice(0, -5)
    return expr.slice(0, -1)
}

/** Returns a user-facing error message, or null if the expression can be sent. */
export function validate(expr: string): string | null {
    if (expr.trim() === '') return 'Enter an expression'

    let depth = 0
    for (const c of expr) {
        if (c === '(') depth++
        else if (c === ')') depth--
        if (depth < 0) return 'Unexpected closing parenthesis' // extra ")"
    }
    if (depth > 0) return 'Missing closing parenthesis'
    return null
}

/** Pretty-prints for the display: * → ×, / → ÷, sqrt → √ */
export function toDisplay(expr: string): string {
    return expr.replace(/\*/g, '×').replace(/\//g, '÷').replace(/sqrt\(/g, '√(')
}

/**
 * Trims floating-point noise (0.30000000000000004 → "0.3").
 *
 * This is a presentation decision: the backend returns the exact float64 and
 * the frontend decides how to show it. toPrecision(12) rounds to 12
 * significant digits (hiding binary rounding error) and Number(...) drops the
 * trailing zeros it leaves behind. Very large or small values keep exponent
 * notation (1e21 → "1e+21").
 */
export function formatResult(n: number): string {
    return Number(n.toPrecision(12)).toString()
}