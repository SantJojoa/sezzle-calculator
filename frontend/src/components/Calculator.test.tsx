import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Calculator } from './Calculator'
import { evaluate } from '../api/calculatorClient'

vi.mock('../api/calculatorClient', async (importOriginal) => ({
    ...(await importOriginal<typeof import('../api/calculatorClient')>()),
    evaluate: vi.fn(),
}))

const evaluateMock = vi.mocked(evaluate)

describe('Calculator accessibility and keyboard handling', () => {
    beforeEach(() => {
        evaluateMock.mockReset()
    })
    afterEach(() => {
        cleanup() // vitest globals are off, so RTL does not auto-clean
        vi.restoreAllMocks()
    })

    it('announces the result in the status live region', async () => {
        evaluateMock.mockResolvedValue(14)
        const user = userEvent.setup()
        render(<Calculator />)

        expect(screen.getByRole('status')).toBeEmptyDOMElement() // mounted before any result

        await user.keyboard('2+3*4{Enter}')

        expect(evaluateMock).toHaveBeenCalledWith('2+3*4')
        expect(await screen.findByRole('status')).toHaveTextContent('14')
    })

    it('activates a focused key with Enter instead of submitting', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        screen.getByRole('button', { name: '7' }).focus()
        await user.keyboard('{Enter}')

        expect(screen.getByTestId('expression')).toHaveTextContent('7')
        expect(evaluateMock).not.toHaveBeenCalled()
    })

    it('still submits with Enter after clicking a key with the mouse', async () => {
        evaluateMock.mockResolvedValue(7)
        const user = userEvent.setup()
        render(<Calculator />)

        await user.click(screen.getByRole('button', { name: '7' }))
        await user.keyboard('{Enter}')

        expect(evaluateMock).toHaveBeenCalledWith('7')
    })

    it('labels the keypad group and the symbol keys', () => {
        render(<Calculator />)

        expect(screen.getByRole('group', { name: 'Calculator keypad' })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'multiply' })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'square root' })).toBeInTheDocument()
    })

    it('shows errors in an alert', async () => {
        const { ApiError } = await import('../api/calculatorClient')
        evaluateMock.mockRejectedValue(new ApiError('DIVISION_BY_ZERO', 'cannot divide by zero', 422))
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('1/0{Enter}')

        expect(await screen.findByRole('alert')).toHaveTextContent('cannot divide by zero')
    })
})
