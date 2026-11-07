import { toDisplay } from '../lib/expression'

interface DisplayProps {
    expression: string
    result: string | null
    error: string | null
    loading?: boolean
}

export function Display({ expression, result, error, loading = false }: DisplayProps) {
    return (
        <div className="rounded-2xl bg-display p-4 text-right text-white short:flex short:flex-1 short:flex-col short:justify-end short:p-2 md:p-6">
            <div
                className="min-h-7 break-all text-base text-muted md:text-xl short:min-h-5 short:text-sm"
                data-testid="expression"
            >
                {expression !== '' && <span className="sr-only">Expression: </span>}
                {toDisplay(expression)}
            </div>

            {/* live region must always stay mounted (even empty): only its content changes */}
            <div
                className="min-h-12 break-all text-4xl font-bold tabular-nums text-white md:text-5xl short:min-h-9 short:text-3xl"
                role="status"
                aria-live="polite"
                aria-atomic="true"
                aria-busy={loading}
                data-testid="result"
            >
                {result !== null && <span className="sr-only">Result: </span>}
                {result}
            </div>

            {/* rendered only when there is an error: role="alert" announces on insertion */}
            {error && (
                <p role="alert" className="mt-2 text-sm text-error">
                    {error}
                </p>
            )}
        </div>
    )
}
