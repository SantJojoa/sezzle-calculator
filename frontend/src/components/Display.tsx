import { toDisplay } from '../lib/expression'

interface DisplayProps {
    expression: string
    result: string | null
    error: string | null
    loading?: boolean
}

export function Display({ expression, result, error, loading = false }: DisplayProps) {
    return (
        <div className="rounded-lg bg-gray-900 p-4 text-right text-white" aria-busy={loading}>
            <div
                className="min-h-8 overflow-x-auto break-all text-xl text-gray-400"
                data-testid="expression"
            >
                {toDisplay(expression)}
            </div>

            {/* live region stays mounted so screen readers announce changes to it */}
            <div
                className="min-h-14 break-all text-4xl font-semibold"
                aria-live="polite"
                data-testid="result"
            >
                {result}
            </div>

            {/* rendered only when there is an error: role="alert" announces on insertion */}
            {error && (
                <p role="alert" className="mt-2 text-sm text-red-400">
                    {error}
                </p>
            )}
        </div>
    )
}
