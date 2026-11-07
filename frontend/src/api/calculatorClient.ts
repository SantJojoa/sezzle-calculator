export class ApiError extends Error {
    readonly code: string
    readonly status?: number

    constructor(code: string, message: string, status?: number) {
        super(message)
        this.name = 'ApiError'
        this.code = code
        this.status = status
    }
}

interface ErrorBody {
    error: { code: string, message: string }
}
interface ResultBody {
    result: number
}

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

export async function evaluate(expression: string): Promise<number> {
    let res: Response
    try {
        res = await fetch(`${API_BASE}/api/v1/evaluate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ expression })
        })
    } catch {
        throw new ApiError('NETWORK_ERROR', 'Failed to connect to server', undefined)
    }
    if (!res.ok) {
        let body: unknown
        try {
            body = await res.json()
        } catch {
            // non-JSON body, e.g. a plain-text 404/405 from the router
            throw new ApiError('UNEXPECTED_RESPONSE', 'Unexpected server response', res.status)
        }
        if (isErrorBody(body)) {
            throw new ApiError(body.error.code, body.error.message, res.status)
        }
        throw new ApiError('UNEXPECTED_RESPONSE', 'Unexpected server response', res.status)
    }

    let body: unknown
    try {
        body = await res.json()
    } catch {
        throw new ApiError('UNEXPECTED_RESPONSE', 'Unexpected server response', res.status)
    }
    if (!isResultBody(body)) {
        throw new ApiError('UNEXPECTED_RESPONSE', 'Unexpected server response', res.status)
    }
    return body.result
}

function isErrorBody(v: unknown): v is ErrorBody {
    if (typeof v !== 'object' || v === null || !('error' in v)) return false
    const e = (v as { error: unknown }).error
    return (
        typeof e === 'object' && e !== null &&
        typeof (e as { code?: unknown }).code === 'string' &&
        typeof (e as { message?: unknown }).message === 'string'
    )
}

function isResultBody(v: unknown): v is ResultBody {
    return typeof v === 'object' && v !== null &&
        typeof (v as { result?: unknown }).result === 'number'
}
