import { toDisplay } from '../lib/expression'

interface DisplayProps {
    expression: string
    result: string | null
    error: string | null
    loading?: boolean
}

export function Display({ expression, result, error, loading = false }: DisplayProps) {
    return (
        <div className="rounded-2xl bg-display p-4 text-right text-white" aria-busy={loading}>
            <div
                className="min-h-7 break-all text-base text-muted"
                data-testid="expression"
            >
                {toDisplay(expression)}
            </div>

            {/* live region stays mounted so screen readers announce changes to it */}
            <div
                className="min-h-12 break-all text-4xl font-bold tabular-nums text-white"
                aria-live="polite"
                data-testid="result"
            >
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
