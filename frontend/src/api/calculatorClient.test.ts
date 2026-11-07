import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, evaluate } from './calculatorClient'

function mockFetch(impl: () => Promise<Response>) {
    const fetchMock = vi.fn(impl)
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
}

async function catchError(promise: Promise<unknown>): Promise<unknown> {
    try {
        await promise
    } catch (e) {
        return e
    }
    throw new Error('expected promise to reject')
}

describe('evaluate', () => {
    afterEach(() => {
        vi.unstubAllGlobals()
        vi.restoreAllMocks()
    })

    it('returns the result and posts the expression as JSON', async () => {
        const fetchMock = mockFetch(async () => Response.json({ result: 14 }))

        await expect(evaluate('2+3*4')).resolves.toBe(14)

        expect(fetchMock).toHaveBeenCalledTimes(1)
        const [url, init] = fetchMock.mock.calls[0] as unknown as [string, RequestInit]
        expect(url).toBe('/api/v1/evaluate')
        expect(init.method).toBe('POST')
        expect(init.headers).toEqual({ 'Content-Type': 'application/json' })
        expect(JSON.parse(init.body as string)).toEqual({ expression: '2+3*4' })
    })

    it('throws ApiError with code, message and status for a backend JSON error', async () => {
        mockFetch(async () =>
            Response.json(
                { error: { code: 'DIVISION_BY_ZERO', message: 'cannot divide by zero' } },
                { status: 422 },
            ),
        )

        const err = await catchError(evaluate('10/0'))

        expect(err).toBeInstanceOf(ApiError)
        expect(err).toMatchObject({
            name: 'ApiError',
            code: 'DIVISION_BY_ZERO',
            message: 'cannot divide by zero',
            status: 422,
        })
    })

    it('throws UNEXPECTED_RESPONSE for a non-JSON error response (text/plain 405)', async () => {
        mockFetch(async () =>
            new Response('Method Not Allowed', {
                status: 405,
                headers: { 'Content-Type': 'text/plain; charset=utf-8' },
            }),
        )

        const err = await catchError(evaluate('1+1'))

        expect(err).toBeInstanceOf(ApiError)
        expect(err).toMatchObject({
            code: 'UNEXPECTED_RESPONSE',
            message: 'Unexpected server response',
            status: 405,
        })
    })

    it('throws NETWORK_ERROR without a status when fetch rejects', async () => {
        mockFetch(async () => {
            throw new TypeError('Failed to fetch')
        })

        const err = await catchError(evaluate('1+1'))

        expect(err).toBeInstanceOf(ApiError)
        expect(err).toMatchObject({
            code: 'NETWORK_ERROR',
            message: 'Failed to connect to server',
        })
        expect((err as ApiError).status).toBeUndefined()
    })
})
