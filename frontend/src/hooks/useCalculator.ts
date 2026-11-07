import { useCallback, useRef, useState } from 'react'
import { ApiError, evaluate } from '../api/calculatorClient'
import { OPERATORS, appendKey, backspace, formatResult, validate } from '../lib/expression'

interface CalculatorState {
    expression: string
    result: string | null
    error: string | null
    loading: boolean
}

const initialState: CalculatorState = { expression: '', result: null, error: null, loading: false }

const isOperator = (key: string) => (OPERATORS as readonly string[]).includes(key)

// A result can seed the next expression only if the backend can parse it back
// (formatResult may produce exponent notation such as "1e+21", which it cannot).
const isContinuable = (result: string) => /^-?\d+(\.\d+)?$/.test(result)

export function useCalculator() {
    const [state, setState] = useState<CalculatorState>(initialState)

    // Identifies the latest edit/request so a slow response for an expression
    // the user has already changed or cleared is ignored.
    const requestId = useRef(0)

    const press = useCallback((key: string) => {
        requestId.current++
        setState((prev) => {
            let expression = prev.expression
            if (prev.result !== null) {
                if (isOperator(key)) {
                    // continue from the result: "14" + "*" -> "14*"
                    if (!isContinuable(prev.result)) return { ...prev, loading: false }
                    expression = prev.result
                } else {
                    // a digit (or anything else) starts a new expression
                    expression = ''
                }
            }
            return {
                expression: appendKey(expression, key),
                result: null,
                error: null,
                loading: false,
            }
        })
    }, [])

    const clear = useCallback(() => {
        requestId.current++
        setState(initialState)
    }, [])

    const deleteLast = useCallback(() => {
        requestId.current++
        setState((prev) => ({
            expression: backspace(prev.expression),
            result: null,
            error: null,
            loading: false,
        }))
    }, [])

    const submit = useCallback(async () => {
        if (state.loading) return

        const { expression } = state
        const message = validate(expression)
        if (message) {
            setState((prev) => ({ ...prev, error: message, result: null }))
            return
        }

        const id = ++requestId.current
        setState((prev) => ({ ...prev, loading: true, error: null, result: null }))
        try {
            const n = await evaluate(expression)
            if (id !== requestId.current) return
            setState((prev) => ({ ...prev, result: formatResult(n), loading: false }))
        } catch (e) {
            if (id !== requestId.current) return
            const error = e instanceof ApiError ? e.message : 'Something went wrong'
            setState((prev) => ({ ...prev, error, loading: false }))
        }
    }, [state])

    return { ...state, press, clear, deleteLast, submit }
}
